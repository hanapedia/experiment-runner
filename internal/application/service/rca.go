package service

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/hanapedia/experiment-runner/internal/application/port"
	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/pkg/utility"
)

// RCAExperimentRunner defines the core service logic.
type RCAExperimentRunner struct {
	config           *domain.ExperimentConfig
	kubernetesClient port.KubernetesClientPort
	chaosExperiment  port.ChaosExperimentsPort
}

// NewExperimentRunner creates new ExperimentRunner instance.
func NewExperimentRunner(config *domain.ExperimentConfig, kubernetesClient port.KubernetesClientPort, chaosExperimentClient port.ChaosExperimentsPort) *RCAExperimentRunner {
	return &RCAExperimentRunner{
		config:           config,
		kubernetesClient: kubernetesClient,
		chaosExperiment:  chaosExperimentClient,
	}
}

// RunExperiments runs the core service logic.
func (runner *RCAExperimentRunner) Run() error {
	// Start LoadGenerator
	runner.config.LoadGeneratorConfig.TotalArrivalRate = strings.Split(runner.config.ArrivalRates, ",")[0]
	runner.config.Duration = (runner.config.RCAConfig.GetDuration() * 2).String()
	runner.config.UpdateNamesWithArrivalRate()
	err := runner.kubernetesClient.CreateLoadGeneratorDeployment(runner.config)
	if err != nil {
		return err
	}
	slog.Info("Started loadgenerator. waiting 1 minute for start up.", "arrival-rate", runner.config.LoadGeneratorConfig.TotalArrivalRate)
	if !runner.config.DryRun {
		time.Sleep(time.Minute)
	}

	// run for normal period
	slog.Info("[Normal Period Start]: Sleeping.", "duration", runner.config.RCAConfig.NormalDuration)
	if !runner.config.DryRun {
		time.Sleep(runner.config.RCAConfig.NormalDuration)
	}

	// query metrics for normal period
	err = runner.kubernetesClient.CreateMetricsProcessorJob(
		runner.config,
		utility.GetTimestampedName(runner.config.ExperimentName+"-normal"),
		utility.GetS3Key(runner.config.MetricsProcessorConfig.S3BucketDir, "normal"),
		runner.config.RCAConfig.NormalDuration,
	)
	if err != nil {
		return err
	}
	slog.Info("[Normal Period End]: Waiting for Injection to start")

	// retrieve deployments
	deployments, err := runner.kubernetesClient.GetDeploymentsWithOutAnnotation(runner.config)
	if err != nil {
		return err
	}
	slog.Info("[Deployment retrieved]: Starting Cycle.", "num-deployment", len(deployments))

	for i, deployment := range deployments {
		slog.Info("[Experiment Start]: Cycle started.", "deployment", deployment.Name)
		slog.Info("[Before Injection]: Sleeping.", "duration", runner.config.RCAConfig.InjectionDuration)
		if !runner.config.DryRun {
			time.Sleep(runner.config.RCAConfig.InjectionDuration)
		}

		err = runner.chaosExperiment.CreateAndApplyNetworkDelay(deployment, runner.config)
		if err != nil {
			return err
		}
		slog.Info("[Injection Period Start]: Injected. Sleeping.", "deployment", deployment.Name, "duration", runner.config.RCAConfig.InjectionDuration)
		if !runner.config.DryRun {
			time.Sleep(runner.config.RCAConfig.InjectionDuration)
		}

		// TODO: must replace with new metrics processor
		slog.Info("[Injection Period End]: Waiting for metrics export to complete",
			"duration",
			runner.config.RCAConfig.InjectionDuration+runner.config.RCAConfig.InjectionDuration/2,
		)
		err = runner.kubernetesClient.CreateMetricsProcessorJob(
			runner.config,
			utility.GetTimestampedName(runner.config.ExperimentName+"-"+deployment.Name),
			utility.GetS3Key(runner.config.MetricsProcessorConfig.S3BucketDir, deployment.Name),
			runner.config.RCAConfig.InjectionDuration+runner.config.RCAConfig.InjectionDuration/2,
		)
		if err != nil {
			return err
		}
		slog.Info(fmt.Sprintf("[Experiment End]: Cycle completed. (%v/%v Done)", i+1, len(deployments)), "deployment", deployment.Name)
		slog.Info("[Draining]: Sleeping for 1 minute.")
		if !runner.config.DryRun {
			time.Sleep(time.Minute)
		}
	}

	slog.Info("Duration complete. deleting loadgenerator deployment")
	err = runner.kubernetesClient.DeleteLoadGeneratorDeployment(runner.config)
	if err != nil {
		return err
	}

	return nil
}

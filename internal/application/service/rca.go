package service

import (
	"log/slog"
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
	deployments, err := runner.kubernetesClient.GetDeploymentsWithOutAnnotation(runner.config)
	if err != nil {
		return err
	}

	for i, deployment := range deployments {
		slog.Info("[Experiment Start]: Cycle started.", "deployment", deployment.Name)
		slog.Info("[Normal Period Start]: Sleeping.", "duration", runner.config.RCAConfig.NormalDuration)
		if !runner.config.DryRun {
			time.Sleep(runner.config.RCAConfig.NormalDuration)
		}
		slog.Info("[Normal Period End]: Waiting for Injection to start")

		err = runner.chaosExperiment.CreateAndApplyNetworkDelay(deployment, runner.config)
		if err != nil {
			return err
		}
		slog.Info("[Injection Period Start]: Injected. Sleeping.", "deployment", deployment.Name, "duration", runner.config.RCAConfig.InjectionDuration)
		if !runner.config.DryRun {
			time.Sleep(runner.config.RCAConfig.InjectionDuration)
		}

		// TODO: must replace with new metrics processor
		slog.Info("[Injection Period End]: Waiting for metrics export to complete")
		err = runner.kubernetesClient.CreateMetricsProcessorJob(
			runner.config,
			utility.GetTimestampedName(runner.config.ExperimentName+"-"+deployment.Name),
			utility.GetS3Key(runner.config.MetricsProcessorConfig.S3BucketDir, deployment.Name),
			runner.config.RCAConfig.GetDuration(),
		)
		if err != nil {
			return err
		}
		slog.Info("[Experiment End]: Cycle completed for '%s'. (%v/%v Done)", deployment.Name, i+1, len(deployments))
		slog.Info("[Draining]: Sleeping for another %s", runner.config.RCAConfig.InjectionDuration)
		if !runner.config.DryRun {
			time.Sleep(runner.config.RCAConfig.InjectionDuration)
		}
	}
	return nil
}

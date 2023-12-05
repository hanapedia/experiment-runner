package service

import (
	"log"
	"log/slog"
	"time"

	"github.com/hanapedia/experiment-runner/internal/application/port"
	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/internal/infrastructure/config"
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
	deployments, err := runner.kubernetesClient.GetDeploymentsWithOutAnnotation(
		runner.config.TargetNamespace,
		config.GetEnvs().RCA_INJECTION_IGNORE_KEY,
		config.GetEnvs().RCA_INJECTION_IGNORE_VALUE,
	)
	if err != nil {
		return err
	}
	for i, deployment := range deployments {
		slog.Info("[INFO]:[Experiment Start]: Cycle started for '%s'", deployment.Name)
		slog.Info("[INFO]:[Normal Period Start]: Sleeping for %s", runner.config.RCAConfig.NormalDuration)
		time.Sleep(runner.config.RCAConfig.NormalDuration)
		slog.Info("[INFO]:[Normal Period End]: Waiting for Injection to start")

		err = runner.chaosExperiment.CreateAndApplyNetworkDelay(deployment)
		if err != nil {
			return err
		}
		slog.Info("[INFO]:[Injection Period Start]: Injected to '%s' Sleeping for %s", deployment.Name, runner.config.RCAConfig.InjectionDuration)
		time.Sleep(runner.config.RCAConfig.InjectionDuration)

		// TODO: must replace with new metrics processor
		slog.Info("[INFO]:[Injection Period End]: Waiting for metrics export to complete")
		err = runner.kubernetesClient.CreateAndApplyJobResource(deployment, runner.config)
		if err != nil {
			return err
		}
		slog.Info("[INFO]:[Experiment End]: Cycle completed for '%s'. (%v/%v Done)", deployment.Name, i+1, len(deployments))
		slog.Info("[INFO]:[Draining]: Sleeping for another %s", runner.config.RCAConfig.InjectionDuration)
		time.Sleep(runner.config.RCAConfig.InjectionDuration)
	}
	return nil
}

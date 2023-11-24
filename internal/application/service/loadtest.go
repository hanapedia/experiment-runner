package service

import (
	"log/slog"
	"time"

	"github.com/hanapedia/experiment-runner/internal/application/port"
	"github.com/hanapedia/experiment-runner/internal/domain"
)

type LoadTestRunner struct {
	kubernetesClient         port.KubernetesClientPort
	metricsQueryConfig       *domain.MetricsProcessorConfig
	loadgeneratorQueryConfig *domain.LoadGeneratorConfig
}

var dryDuration = 1 * time.Minute

func NewLoadTestRunner(kc port.KubernetesClientPort, mqc *domain.MetricsProcessorConfig, lgc *domain.LoadGeneratorConfig) *LoadTestRunner {
	return &LoadTestRunner{
		kubernetesClient:         kc,
		metricsQueryConfig:       mqc,
		loadgeneratorQueryConfig: lgc,
	}
}

func (runner *LoadTestRunner) Run() error {
	err := runner.kubernetesClient.CreateLoadGeneratorDeployment(runner.loadgeneratorQueryConfig)
	if err != nil {
		return err
	}
	slog.Info("[Experiement Started]: sleeping for dry duration", "duration", dryDuration)
	time.Sleep(dryDuration)

	slog.Info("[Experiement Started]: sleeping for load test duration", "duration", runner.metricsQueryConfig.GetDuration())
	time.Sleep(runner.metricsQueryConfig.GetDuration())

	err = runner.kubernetesClient.CreateMetricsProcessorJob(runner.metricsQueryConfig)
	if err != nil {
		return nil
	}

	slog.Info("[Experiement Ended]")

	return nil
}

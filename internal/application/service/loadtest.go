package service

import (
	"log/slog"
	"strings"
	"time"

	"github.com/hanapedia/experiment-runner/internal/application/port"
	"github.com/hanapedia/experiment-runner/internal/domain"
)

type LoadTestRunner struct {
	kubernetesClient port.KubernetesClientPort
	config           *domain.ExperimentConfig
}

var dryDuration = 1 * time.Minute

func NewLoadTestRunner(kc port.KubernetesClientPort, config *domain.ExperimentConfig) *LoadTestRunner {
	return &LoadTestRunner{
		kubernetesClient: kc,
	}
}

func (runner *LoadTestRunner) Run() error {
	rpss := strings.Split(runner.config.ArrivalRates, ",")

	for _, rps := range rpss {
		// update vars
		runner.config.LoadGeneratorConfig.TotalArrivalRate = rps
		runner.config.UpdateNamesWithArrivalRate()

		err := runner.kubernetesClient.CreateLoadGeneratorDeployment(runner.config)
		if err != nil {
			return err
		}
		slog.Info("Started loadgenerator", "arrival-rate", runner.config.LoadGeneratorConfig.TotalArrivalRate)

		slog.Info("[Experiement Started]: sleeping for dry duration", "duration", dryDuration)
		time.Sleep(dryDuration)

		slog.Info("[Experiement Started]: sleeping for load test duration", "duration", runner.config.GetDuration())
		time.Sleep(runner.config.GetDuration())

		err = runner.kubernetesClient.CreateMetricsProcessorJob(runner.config)
		if err != nil {
			return nil
		}
		slog.Info("Started metrics processor", "arrival-rate", runner.config.LoadGeneratorConfig.TotalArrivalRate)
	}
	slog.Info("[Experiement Ended]")

	return nil
}

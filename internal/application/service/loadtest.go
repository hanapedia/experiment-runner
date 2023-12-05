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

type LoadTestRunner struct {
	kubernetesClient port.KubernetesClientPort
	config           *domain.ExperimentConfig
}

var dryDuration = 1 * time.Minute

func NewLoadTestRunner(kc port.KubernetesClientPort, config *domain.ExperimentConfig) *LoadTestRunner {
	return &LoadTestRunner{
		kubernetesClient: kc,
		config:           config,
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
		if !runner.config.DryRun {
			time.Sleep(dryDuration)
		}

		slog.Info("[Experiement Started]: sleeping for load test duration", "duration", runner.config.GetDuration())
		if !runner.config.DryRun {
			time.Sleep(runner.config.GetDuration())
		}

		slog.Info("Duration complete. deleting loadgenerator deployment")
		err = runner.kubernetesClient.DeleteLoadGeneratorDeployment(runner.config)
		if err != nil {
			return err
		}

		err = runner.kubernetesClient.CreateMetricsProcessorJob(
			runner.config,
			utility.GetTimestampedName(fmt.Sprintf("%s-%s", runner.config.ExperimentName, runner.config.K6TestName)),
			utility.GetS3Key(runner.config.MetricsProcessorConfig.S3BucketDir, runner.config.K6TestName),
			runner.config.GetDuration(),
		)
		if err != nil {
			return err
		}
		slog.Info("Started metrics processor", "arrival-rate", runner.config.LoadGeneratorConfig.TotalArrivalRate)
	}
	slog.Info("[Experiement Ended]")

	return nil
}

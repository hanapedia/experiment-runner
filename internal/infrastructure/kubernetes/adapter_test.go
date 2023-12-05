package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/pkg/utility"
)

func TestCreateMetricsProcessorJob(t *testing.T) {
	// Setup: Create a mock KubernetesAdapter and a MetricsProcessorConfig
	adapter := &KubernetesAdapter{}
	config := &domain.ExperimentConfig{
		// common
		ExperimentName:      "test",
		TargetNamespace:     "emulation",
		ExperimentNamespace: "experiment",

		MetricsProcessorConfig: &domain.MetricsProcessorConfig{
			// metrics processor
			ConfigMapName: "metrics-processor",
			ImageName:     "metrics-processor",
			ImageTag:      "v1.0.0",
			S3BucketDir:   "test",
		},
		LoadGeneratorConfig: &domain.LoadGeneratorConfig{
			// metrics processor
			ConfigMapName:    "test-lg",
			ImageName:        "grafana/k6",
			ImageTag:         "0.47.0",
			TotalArrivalRate: "100",
			FrontendAddr:     "frontend:80",
			IndexRoute:       "/",
			K6PromWriteUrl:   "http://prometheus-kube-prometheus-prometheus.monitoring.svc.cluster.local:9090/api/v1/write",
			K6TrendStats:     "p(95),p(99),avg",
		},
	}

	// Test 1: Normal execution
	t.Run("NormalExecution", func(t *testing.T) {
		// Call CreateMetricsProcessorJob
		_ = adapter.CreateMetricsProcessorJob(
			config,
			utility.GetTimestampedName(fmt.Sprintf("%s-%s", config.ExperimentName, config.K6TestName)),
			utility.GetS3Key(config.MetricsProcessorConfig.S3BucketDir, config.K6TestName),
			config.GetDuration(),
		)
	})

	// Test 2: Normal execution
	t.Run("NormalExecution", func(t *testing.T) {
		_ = adapter.CreateLoadGeneratorDeployment(config)
	})

}

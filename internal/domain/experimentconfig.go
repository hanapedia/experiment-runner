package domain

import (
	"fmt"
	"log/slog"
	"time"

	corev1 "k8s.io/api/core/v1"
)

type ExperimentConfig struct {
	// common
	ExperimentName      string
	TargetNamespace     string
	ExperimentNamespace string
	K6TestName          string
	Duration            string
	ArrivalRates        string

	RCAConfig              *RCAExperimentConfig
	MetricsProcessorConfig *MetricsProcessorConfig
	LoadGeneratorConfig    *LoadGeneratorConfig

	// testing
	DryRun bool
}

func (config ExperimentConfig) GetDuration() time.Duration {
	duration, err := time.ParseDuration(config.Duration)
	if err != nil {
		slog.Error("Failed to parse duration string. defaulting to 1 minute.", "err", err)
		return time.Minute
	}
	return duration
}

func (config *ExperimentConfig) UpdateNamesWithArrivalRate() {
	config.K6TestName = config.GetNameWithArrivalRate()
	config.MetricsProcessorConfig.S3BucketDir = config.GetNameWithArrivalRate()
}

func (config ExperimentConfig) GetNameWithArrivalRate() string {
	return fmt.Sprintf("%s-%s", config.ExperimentName, config.LoadGeneratorConfig.TotalArrivalRate)
}

func (config ExperimentConfig) GetLoadGeneratorEnv() []corev1.EnvVar {
	return []corev1.EnvVar{
		{Name: "TEST_NAME", Value: config.K6TestName},
		{Name: "TOTAL_ARRIVAL_RATE", Value: config.LoadGeneratorConfig.TotalArrivalRate},
		{Name: "FRONTEND_ADDR", Value: config.LoadGeneratorConfig.FrontendAddr},
		{Name: "INDEX_ROUTE", Value: config.LoadGeneratorConfig.IndexRoute},
		{Name: "K6_PROMETHEUS_RW_SERVER_URL", Value: config.LoadGeneratorConfig.K6PromWriteUrl},
		{Name: "K6_PROMETHEUS_RW_TREND_STATS", Value: config.LoadGeneratorConfig.K6TrendStats},
	}
}

package config

import "github.com/hanapedia/experiment-runner/internal/domain"

func NewExperimentConfig() *domain.ExperimentConfig {
	return &domain.ExperimentConfig{
		ExperimentName:         GetEnvs().EXPERIMENT_NAME,
		ExperimentNamespace:    GetEnvs().EXPERIMENT_NAMESPACE,
		TargetNamespace:        GetEnvs().TARGET_NAMESPACE,
		K6TestName:             GetEnvs().K6_TEST_NAME,
		Duration:               GetEnvs().DURATION,
		ArrivalRates:           GetEnvs().ARRIVAL_RATES,
		RCAConfig:              NewRCAExperimentConfig(),
		MetricsProcessorConfig: NewMetricsProcessorConfig(),
		LoadGeneratorConfig:    NewLoadGeneratorConfig(),
	}
}

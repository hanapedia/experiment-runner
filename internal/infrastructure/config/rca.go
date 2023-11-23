package config

import (
	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/pkg/utility"
)

func NewRCAExperimentConfig() *domain.RCAExperimentConfig {
	return &domain.RCAExperimentConfig{
		Name:                      GetEnvs().EXPERIMENT_NAME,
		TargetNamespace:           GetEnvs().TARGET_NAMESPACE,
		ExperimentNamespace:       GetEnvs().EXPERIMENT_NAMESPACE,
		MetricsQueryConfigMapName: GetEnvs().METRICS_QUERY_CONFIG_MAP_NAME,
		MetricsQueryImageName:     GetEnvs().METRICS_QUERY_IMAGE,
		NormalDuration: utility.ParseDurationWithDefault(
			GetEnvs().RCA_NORMAL_DURATION,
			domain.DefaultRCAExperimentConfig.NormalDuration,
		),
		InjectionDuration: utility.ParseDurationWithDefault(
			GetEnvs().RCA_INJECTION_DURATION,
			domain.DefaultRCAExperimentConfig.InjectionDuration,
		),
		Latency: utility.ParseDurationWithDefault(
			GetEnvs().RCA_LATENCY,
			domain.DefaultRCAExperimentConfig.Latency,
		),
		Jitter: utility.ParseDurationWithDefault(
			GetEnvs().RCA_JITTER,
			domain.DefaultRCAExperimentConfig.Jitter,
		),
	}
}

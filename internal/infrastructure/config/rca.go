package config

import (
	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/pkg/utility"
)

func NewRCAExperimentConfig() *domain.RCAExperimentConfig {
	return &domain.RCAExperimentConfig{
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

package domain

import (
	"time"
)

type RCAExperimentConfig struct {
	// NormalDuration is the duration without injection
	NormalDuration time.Duration

	// InjectionDuration is the duration of injection
	InjectionDuration time.Duration

	// Latency is the amount of network delay to inject
	Latency time.Duration

	// Jitter is the variance in the amount of network delay injected
	Jitter time.Duration

	RcaInjectionIgnoreKey   string
	RcaInjectionIgnoreValue string
}

var DefaultRCAExperimentConfig = RCAExperimentConfig{
	NormalDuration:    5 * time.Minute,
	InjectionDuration: 1 * time.Minute,
	Latency:           15 * time.Millisecond,
	Jitter:            5 * time.Millisecond,
}

func (config RCAExperimentConfig) GetDuration() time.Duration {
	return config.InjectionDuration + config.NormalDuration
}

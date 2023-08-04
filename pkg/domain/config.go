package domain

import (
	"strconv"
	"time"
)

type ExperimentConfig struct {
	// Name is the name of the experiment
	Name string

	// TargetNamespace is the namespace that application is running in
	TargetNamespace string

	// NormalDuration is the duration without injection
	NormalDuration time.Duration

	// InjectionDuration is the duration of injection
	InjectionDuration time.Duration

	// Latency is the amount of network delay to inject
	Latency time.Duration

	// Jitter is the variance in the amount of network delay injected
	Jitter time.Duration
}

func (config ExperimentConfig) GetDuration() string {
	seconds := int((config.InjectionDuration + config.NormalDuration).Seconds())
	return strconv.Itoa(seconds)
}

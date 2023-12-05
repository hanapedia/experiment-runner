package domain

import (
	"strconv"
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
}

var DefaultRCAExperimentConfig = RCAExperimentConfig{
	NormalDuration:    5 * time.Minute,
	InjectionDuration: 1 * time.Minute,
	Latency:           15 * time.Millisecond,
	Jitter:            5 * time.Millisecond,
}

func (config RCAExperimentConfig) GetDuration() string {
	seconds := int((config.InjectionDuration + config.NormalDuration).Seconds())
	return strconv.Itoa(seconds)
}

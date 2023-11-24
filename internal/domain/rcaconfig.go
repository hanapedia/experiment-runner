package domain

import (
	"fmt"
	"strconv"
	"time"
)

type RCAExperimentConfig struct {
	// Name is the name of the experiment
	Name string

	// TargetNamespace is the namespace that application is running in
	TargetNamespace string

	// TargetNamespace is the namespace that application is running in
	ExperimentNamespace string

	// MetricsProcessorConfigMapName is the name of config map that will be fed to the Job created
	MetricsProcessorConfigMapName string

	// MetricsProcessorImageName is the tag of the batch job container image
	MetricsProcessorImageName string

	// MetricsProcessorImageName is the tag of the batch job container image
	MetricsProcessorImageTag string

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

func (config RCAExperimentConfig) GetMetricsProcessorImageName() string {
	return fmt.Sprintf("%s:%s", config.MetricsProcessorImageName, config.MetricsProcessorImageTag)
}

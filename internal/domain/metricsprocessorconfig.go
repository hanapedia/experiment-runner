package domain

import (
	"fmt"
	"log/slog"
	"time"
)

type MetricsProcessorConfig struct {
	// common
	ExperimentName      string
	TargetNamespace     string
	ExperimentNamespace string

	// metrics processor
	ConfigMapName string
	ImageName     string
	ImageTag      string
	S3BucketDir   string
	K6TestName    string
	Duration      string
}

func (config MetricsProcessorConfig) GetMetricsProcessorImageName() string {
	return fmt.Sprintf("%s:%s", config.ImageName, config.ImageTag)
}

func (config MetricsProcessorConfig) GetDuration() time.Duration {
	duration, err := time.ParseDuration(config.Duration)
	if err != nil {
		slog.Error("Failed to parse duration string. defaulting to 1 minute.", "err", err)
		return time.Minute
	}
	return duration
}

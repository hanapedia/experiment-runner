package domain

import (
	"fmt"
)

type MetricsProcessorConfig struct {
	// metrics processor
	ConfigMapName string
	ImageName     string
	ImageTag      string
	S3BucketDir   string
}

func (config MetricsProcessorConfig) GetImageName() string {
	return fmt.Sprintf("%s:%s", config.ImageName, config.ImageTag)
}

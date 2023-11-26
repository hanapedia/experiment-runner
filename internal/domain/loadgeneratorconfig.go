package domain

import (
	"fmt"
)

type LoadGeneratorConfig struct {
	// metrics processor
	ConfigMapName    string
	ImageName        string
	ImageTag         string
	TotalArrivalRate string
	FrontendAddr     string
	IndexRoute       string
	K6PromWriteUrl   string
	K6TrendStats     string
}

func (config LoadGeneratorConfig) GetImageName() string {
	return fmt.Sprintf("%s:%s", config.ImageName, config.ImageTag)
}


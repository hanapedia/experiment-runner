package domain

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type LoadGeneratorConfig struct {
	// common
	ExperimentName      string
	TargetNamespace     string
	ExperimentNamespace string

	// metrics processor
	ConfigMapName    string
	ImageName        string
	ImageTag         string
	K6TestName       string
	Duration         string
	TotalArrivalRate string
	FrontendAddr     string
	IndexRoute       string
	K6PromWriteUrl   string
	K6TrendStats     string
}

func (config LoadGeneratorConfig) GetLoadGeneratorImageName() string {
	return fmt.Sprintf("%s:%s", config.ImageName, config.ImageTag)
}

func (config LoadGeneratorConfig) MapToEnv() []corev1.EnvVar {
	return []corev1.EnvVar{
		{Name: "TOTAL_ARRIVAL_RATE", Value: config.TotalArrivalRate},
		{Name: "DURATION", Value: config.Duration},
		{Name: "FRONTEND_ADDR", Value: config.FrontendAddr},
		{Name: "INDEX_ROUTE", Value: config.IndexRoute},
		{Name: "TEST_NAME", Value: config.K6TestName},
		{Name: "K6_PROMETHEUS_RW_SERVER_URL", Value: config.K6PromWriteUrl},
		{Name: "K6_PROMETHEUS_RW_TREND_STATS", Value: config.K6TrendStats},
	}

}

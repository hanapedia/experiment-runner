package config

import "github.com/hanapedia/experiment-runner/internal/domain"

func NewLoadGeneratorConfig() *domain.LoadGeneratorConfig {
	return &domain.LoadGeneratorConfig{
		ExperimentName:      GetEnvs().EXPERIMENT_NAME,
		TargetNamespace:     GetEnvs().TARGET_NAMESPACE,
		ExperimentNamespace: GetEnvs().EXPERIMENT_NAME,
		K6TestName:          GetEnvs().K6_TEST_NAME,
		Duration:            GetEnvs().DURATION,
		ConfigMapName:       GetEnvs().LG_CONFIG_MAP_NAME,
		ImageName:           GetEnvs().LG_IMAGE,
		ImageTag:            GetEnvs().LG_IMAGE_TAG,
		TotalArrivalRate:    GetEnvs().LG_TOTAL_ARRIVAL_RATE,
		FrontendAddr:        GetEnvs().LG_FRONTEND_ADDR,
		IndexRoute:          GetEnvs().LG_INDEX_ROUTE,
		K6PromWriteUrl:      GetEnvs().LG_K6_PROMETHEUS_RW_SERVER_URL,
		K6TrendStats:        GetEnvs().LG_K6_PROMETHEUS_RW_TREND_STATS,
	}
}

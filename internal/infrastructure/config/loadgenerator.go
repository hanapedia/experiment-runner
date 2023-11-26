package config

import "github.com/hanapedia/experiment-runner/internal/domain"

func NewLoadGeneratorConfig() *domain.LoadGeneratorConfig {
	return &domain.LoadGeneratorConfig{
		ConfigMapName:       GetEnvs().LG_CONFIG_MAP_NAME,
		ImageName:           GetEnvs().LG_IMAGE,
		ImageTag:            GetEnvs().LG_IMAGE_TAG,
		FrontendAddr:        GetEnvs().LG_FRONTEND_ADDR,
		IndexRoute:          GetEnvs().LG_INDEX_ROUTE,
		K6PromWriteUrl:      GetEnvs().LG_K6_PROMETHEUS_RW_SERVER_URL,
		K6TrendStats:        GetEnvs().LG_K6_PROMETHEUS_RW_TREND_STATS,
	}
}

package config

import "github.com/hanapedia/experiment-runner/internal/domain"

func NewMetricsProcessorConfig() *domain.MetricsProcessorConfig {
	return &domain.MetricsProcessorConfig{
		ConfigMapName:       GetEnvs().METRICS_PROCESSOR_CONFIG_MAP_NAME,
		ImageName:           GetEnvs().METRICS_PROCESSOR_IMAGE,
		ImageTag:            GetEnvs().METRICS_PROCESSOR_IMAGE_TAG,
		S3BucketDir:         GetEnvs().METRICS_PROCESSOR_S3_BUCKET_DIR,
	}
}

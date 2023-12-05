package port

import (
	"time"

	"github.com/hanapedia/experiment-runner/internal/domain"
)

// KubernetesClientPort defines the interface for interactions with the Kubernetes API
type KubernetesClientPort interface {
	// GetDeploymentsWithAnnotation retrieves deployments in the given namespace with the specified annotation
	GetDeploymentsWithOutAnnotation(config *domain.ExperimentConfig) ([]domain.Deployment, error)

	// CreateMetricsProcessorJob creates and starts a job to process metrics
	CreateMetricsProcessorJob(config *domain.ExperimentConfig, name, bucketDir string, duration time.Duration) error

	// CreateLoadGeneratorDeployment creates deployment for load generator pod
	CreateLoadGeneratorDeployment(config *domain.ExperimentConfig) error

	// DeleteLoadGeneratorDeployment delete deployment for load generator pod
	DeleteLoadGeneratorDeployment(config *domain.ExperimentConfig) error
}

// ChaosExperimentsPort defines the interface for interactions with chaos experiment tools
type ChaosExperimentsPort interface {
	// CreateAndApplyNetworkDelay creates a new chaos resource for the specified deployment and apply it.
	CreateAndApplyNetworkDelay(deployment domain.Deployment, config *domain.ExperimentConfig) error
}

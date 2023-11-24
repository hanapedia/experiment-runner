package port

import "github.com/hanapedia/experiment-runner/internal/domain"

// KubernetesClientPort defines the interface for interactions with the Kubernetes API
type KubernetesClientPort interface {
	// GetDeploymentsWithAnnotation retrieves deployments in the given namespace with the specified annotation
	GetDeploymentsWithOutAnnotation(namespace string, annotationKey string, annotationValue string) ([]domain.Deployment, error)
	// CreateAndApplyJobResource creates a new job resource and applies it for the specified deployment
	CreateAndApplyJobResource(deployment domain.Deployment, config *domain.RCAExperimentConfig) error

	// CreateMetricsProcessorJob creates and starts a job to process metrics
	CreateMetricsProcessorJob(config *domain.MetricsProcessorConfig) error

	// CreateLoadGeneratorDeployment creates deployment for load generator pod
	CreateLoadGeneratorDeployment(config *domain.LoadGeneratorConfig) error
}

// ChaosExperimentsPort defines the interface for interactions with chaos experiment tools
type ChaosExperimentsPort interface {
	// CreateAndApplyNetworkDelay creates a new chaos resource for the specified deployment and apply it.
	CreateAndApplyNetworkDelay(deployment domain.Deployment) error
}

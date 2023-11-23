package kubernetes

import (
	"github.com/hanapedia/experiment-runner/internal/application/port"
	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/pkg/utility"
	"k8s.io/client-go/kubernetes"
)

type KubernetesAdapter struct {
	client *KubernetesClient
	config *domain.RCAExperimentConfig
}

func NewKubernetesAdapter(clientset *kubernetes.Clientset, config *domain.RCAExperimentConfig) port.KubernetesClientPort {
	return &KubernetesAdapter{
		client: NewKubernetesClient(clientset),
		config: config,
	}
}

// GetDeploymentsWithOutAnnotation converts retrieved kubernetes api deployments to domain deployments.
func (adapter *KubernetesAdapter) GetDeploymentsWithOutAnnotation(namespace string, annotationKey string, annotationValue string) ([]domain.Deployment, error) {
	deployments, err := adapter.client.GetDeploymentsWithOutAnnotation(namespace, annotationKey, annotationValue)
	if err != nil {
		return nil, err
	}

	var domainDeployments []domain.Deployment
	for _, deployment := range deployments {
		domainDeployments = append(domainDeployments, domain.Deployment{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		})
	}

	return domainDeployments, nil
}

// CreateAndApplyJobResource converts domain job to kubernetes api job type and create Job.
func (adapter *KubernetesAdapter) CreateAndApplyJobResource(deployment domain.Deployment) error {
	job := ConstructJob(JobArgs{
		Name:            utility.GetTimestampedName(adapter.config.Name, deployment.Name),
		S3Key:           utility.GetS3Key(adapter.config.Name, deployment.Name),
		TargetNamespace: adapter.config.TargetNamespace,
		ConfigMapName:   adapter.config.MetricsQueryConfigMapName,
		JobImageName:    adapter.config.GetMetricsQueryImageName(),
		Duration:        adapter.config.GetDuration(),
	},
	)
	err := adapter.client.ApplyJobResource(job, adapter.config.ExperimentNamespace)
	if err != nil {
		return err
	}
	return nil
}

package kubernetes

import (
	"github.com/hanapedia/rca-experiment-runner/pkg/application/port"
	"github.com/hanapedia/rca-experiment-runner/pkg/domain"
	"github.com/hanapedia/rca-experiment-runner/util"
	"k8s.io/client-go/kubernetes"
)

type KubernetesAdapter struct {
	client *KubernetesClient
	config *domain.ExperimentConfig
}

func NewKubernetesAdapter(clientset *kubernetes.Clientset, config *domain.ExperimentConfig) port.KubernetesClientPort {
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
	job := ConstructJob(util.GetTimestampedName(deployment.Name), adapter.config.GetDuration())
	err := adapter.client.ApplyJobResource(job)
	if err != nil {
		return err
	}
	return nil
}

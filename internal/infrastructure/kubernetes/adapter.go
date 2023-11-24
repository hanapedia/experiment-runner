package kubernetes

import (
	"fmt"

	"github.com/hanapedia/experiment-runner/internal/application/port"
	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/pkg/utility"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesAdapter struct {
	client *KubernetesClient
}

func NewKubernetesAdapter(kubeConfig *rest.Config) port.KubernetesClientPort {
	// prepare kube client for kubernetes API
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(err.Error())
	}
	return &KubernetesAdapter{
		client: NewKubernetesClient(clientset),
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
func (adapter *KubernetesAdapter) CreateAndApplyJobResource(deployment domain.Deployment, config *domain.RCAExperimentConfig) error {
	job := ConstructJob(JobArgs{
		Name:            utility.GetTimestampedName(config.Name + "-" + deployment.Name),
		S3BucketDir:     utility.GetS3Key(config.Name, deployment.Name),
		TargetNamespace: config.TargetNamespace,
		ConfigMapName:   config.MetricsProcessorConfigMapName,
		JobImageName:    config.GetMetricsProcessorImageName(),
		Duration:        config.GetDuration(),
	},
	)
	err := adapter.client.ApplyJobResource(job, config.ExperimentNamespace)
	if err != nil {
		return err
	}
	return nil
}

func (adapter *KubernetesAdapter) CreateMetricsProcessorJob(config *domain.MetricsProcessorConfig) error {
	job := ConstructJob(JobArgs{
		Name:            utility.GetTimestampedName(config.ExperimentName),
		S3BucketDir:     config.S3BucketDir,
		TargetNamespace: config.TargetNamespace,
		ConfigMapName:   config.ConfigMapName,
		JobImageName:    config.GetMetricsProcessorImageName(),
		Duration:        config.Duration,
	},
	)
	err := adapter.client.ApplyJobResource(job, config.ExperimentNamespace)
	if err != nil {
		return err
	}
	return nil
}

func (adapter *KubernetesAdapter) CreateLoadGeneratorDeployment(config *domain.LoadGeneratorConfig) error {
	deployment := factory.NewDeployment(&factory.DeploymentArgs{
		Name:     fmt.Sprintf("%s-lg", config.ExperimentName),
		Image:    config.GetLoadGeneratorImageName(),
		Replicas: LG_REPLICAS,
		Resource: &corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    LG_CPU_LIMIT,
				corev1.ResourceMemory: LG_MEMORY_LIMIT,
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    LG_CPU_REQUEST,
				corev1.ResourceMemory: LG_MEMORY_REQUEST,
			},
		},
		VolumeMounts: map[string]string{
			"config": "/scripts/",
		},
		ConfigVolume: &factory.ConfigMapVolumeArgs{
			Name: config.ConfigMapName,
			Items: map[string]string{
				"script.js": "script.js",
			},
		},
		Envs: config.MapToEnv(),
	})

	err := adapter.client.CreateDeployment(&deployment, config.TargetNamespace)
	if err != nil {
		return err
	}
	return nil

}

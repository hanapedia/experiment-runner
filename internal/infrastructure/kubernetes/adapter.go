package kubernetes

import (
	"fmt"
	"time"

	"github.com/hanapedia/experiment-runner/internal/application/port"
	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/pkg/file"
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
func (adapter *KubernetesAdapter) GetDeploymentsWithOutAnnotation(config *domain.ExperimentConfig) ([]domain.Deployment, error) {
	if config.DryRun {
		return []domain.Deployment{
			{Name: "test", Namespace: "test"},
		}, nil
	}
	deployments, err := adapter.client.GetDeploymentsWithOutAnnotation(
		config.TargetNamespace,
		config.RCAConfig.RcaInjectionIgnoreKey,
		config.RCAConfig.RcaInjectionIgnoreValue,
	)
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

func (adapter *KubernetesAdapter) CreateMetricsProcessorJob(config *domain.ExperimentConfig, name, bucketDir string, duration time.Duration) error {
	job := ConstructJob(JobArgs{
		Name:            name,
		S3BucketDir:     bucketDir,
		K6TestName:      config.K6TestName,
		TargetNamespace: config.TargetNamespace,
		ConfigMapName:   config.MetricsProcessorConfig.ConfigMapName,
		JobImageName:    config.MetricsProcessorConfig.GetImageName(),
		Duration:        duration.String(),
	})
	if config.DryRun {
		file.WriteKubernetesManifest(job, fmt.Sprintf("%s-job.yaml", config.K6TestName))
		return nil
	}

	err := adapter.client.ApplyJobResource(job, config.ExperimentNamespace)
	if err != nil {
		return err
	}
	return nil
}

func (adapter *KubernetesAdapter) CreateLoadGeneratorDeployment(config *domain.ExperimentConfig) error {
	deployment := factory.NewDeployment(&factory.DeploymentArgs{
		Name:     fmt.Sprintf("%s-lg", config.ExperimentName),
		Image:    config.LoadGeneratorConfig.GetImageName(),
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
			Name: config.LoadGeneratorConfig.ConfigMapName,
			Items: map[string]string{
				"script.js": "script.js",
			},
		},
		Envs: config.GetLoadGeneratorEnv(),
	})

	deployment.Spec.Template.Spec.Containers[0].Command = []string{
		"k6",
		"run",
		"-o",
		"experimental-prometheus-rw",
		"/scripts/script.js",
	}

	if config.DryRun {
		file.WriteKubernetesManifest(deployment, fmt.Sprintf("%s-deployment.yaml", config.K6TestName))
		return nil
	}

	err := adapter.client.CreateDeployment(&deployment, config.ExperimentNamespace)
	if err != nil {
		return err
	}
	return nil
}

func (adapter *KubernetesAdapter) DeleteLoadGeneratorDeployment(config *domain.ExperimentConfig) error {
	if config.DryRun {
		return nil
	}
	err := adapter.client.DeleteDeployment(fmt.Sprintf("%s-lg", config.ExperimentName), config.ExperimentNamespace)
	if err != nil {
		return err
	}

	return nil
}

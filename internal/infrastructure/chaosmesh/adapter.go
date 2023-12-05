package chaosmesh

import (
	"fmt"
	"strings"

	"github.com/hanapedia/experiment-runner/internal/application/port"
	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/pkg/file"
	"github.com/hanapedia/experiment-runner/pkg/utility"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ChaosMeshAdapter struct {
	client *ChaosMeshClient
}

func NewChaosMeshAdapter(kubeConfig *rest.Config) port.ChaosExperimentsPort {
	// Prepare kube dynamic config for chaos mesh resource
	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		panic(err.Error())
	}
	return &ChaosMeshAdapter{
		client: NewChaosMeshClient(dynamicClient),
	}
}

func (adapter *ChaosMeshAdapter) CreateAndApplyNetworkDelay(deployment domain.Deployment, config *domain.ExperimentConfig) error {
	networkDelay := ConstructNetworkChaos(&NetworkChaosArgs{
		Name:            utility.GetTimestampedName(config.ExperimentName + "-" + deployment.Name),
		TargetNamespace: deployment.Namespace,
		Selector:        map[string]string{"app": deployment.Name},
		Duration:        config.RCAConfig.InjectionDuration.String(),
		Latency:         config.RCAConfig.Latency.String(),
		Jitter:          config.RCAConfig.Jitter.String(),
	})
	if config.DryRun {
		file.WriteKubernetesManifest(deployment, fmt.Sprintf("%s-network-chaos.yaml", deployment.Name))
		return nil
	}
	err := adapter.client.ApplyNetworkDelay(networkDelay)
	if err != nil {
		return err
	}
	return nil
}

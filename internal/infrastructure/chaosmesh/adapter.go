package chaosmesh

import (
	"github.com/hanapedia/experiment-runner/internal/application/port"
	"github.com/hanapedia/experiment-runner/internal/domain"
	"github.com/hanapedia/experiment-runner/pkg/utility"
	"k8s.io/client-go/dynamic"
)

type ChaosMeshAdapter struct {
	client *ChaosMeshClient
	config *domain.ExperimentConfig
}

func NewChaosMeshAdapter(dynamicClient dynamic.Interface, config *domain.ExperimentConfig) port.ChaosExperimentsPort {
	return &ChaosMeshAdapter{
		client: NewChaosMeshClient(dynamicClient),
		config: config,
	}
}

func (adapter *ChaosMeshAdapter) CreateAndApplyNetworkDelay(deployment domain.Deployment) error {
	networkDelay := ConstructNetworkChaos(&NetworkChaosArgs{
		Name:            utility.GetTimestampedName(adapter.config.Name, deployment.Name),
		TargetNamespace: deployment.Namespace,
		Selector:        map[string]string{"app": deployment.Name},
		Duration:        adapter.config.InjectionDuration.String(),
		Latency:         adapter.config.Latency.String(),
		Jitter:          adapter.config.Jitter.String(),
	})
	err := adapter.client.ApplyNetworkDelay(networkDelay)
	if err != nil {
		return err
	}
	return nil
}

package chaosmesh

import (
	"github.com/hanapedia/rca-experiment-runner/pkg/application/port"
	"github.com/hanapedia/rca-experiment-runner/pkg/domain"
	"github.com/hanapedia/rca-experiment-runner/util"
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
		Name:            util.GetTimestampedName(deployment.Name),
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

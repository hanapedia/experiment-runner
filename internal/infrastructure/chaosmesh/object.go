package chaosmesh

import (
	"fmt"

	"github.com/hanapedia/experiment-runner/internal/infrastructure/crd/chaosmesh"
	"github.com/hanapedia/experiment-runner/internal/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetworkChaosArgs struct {
	Name            string
	TargetNamespace string
	Selector        map[string]string
	Duration        string
	Latency         string
	Jitter          string
}

// ConstructNetworkChaos create type network chaos kubernetes custom resource objects for chaos-mesh.
func ConstructNetworkChaos(args *NetworkChaosArgs) *chaosmesh.NetworkChaos {
	return &chaosmesh.NetworkChaos{
		TypeMeta: metav1.TypeMeta{
			Kind:       ChaosMeshNetworkChaosResource,
			APIVersion: fmt.Sprintf("%s/%s", ChaosMeshGroup, ChaosMeshVersion),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      args.Name,
			Namespace: constants.ChaosExperimentNamespace,
		},
		Spec: ConstructNetworkChaosSpec(args),
	}
}

// ConstructNetworkChaosSpec create type namespace kubernetes objects.
func ConstructNetworkChaosSpec(args *NetworkChaosArgs) chaosmesh.NetworkChaosSpec {
	return chaosmesh.NetworkChaosSpec{
		PodSelector: chaosmesh.PodSelector{
			Selector: chaosmesh.PodSelectorSpec{
				GenericSelectorSpec: chaosmesh.GenericSelectorSpec{
					Namespaces:     []string{args.TargetNamespace},
					LabelSelectors: args.Selector,
				},
			},
			Mode: chaosmesh.AllMode,
		},
		Duration:  &args.Duration,
		Action:    chaosmesh.DelayAction,
		Direction: chaosmesh.Both,
		TcParameter: chaosmesh.TcParameter{
			Delay: &chaosmesh.DelaySpec{
				Latency: args.Latency,
				Jitter:  args.Jitter,
			},
		},
		Target: &chaosmesh.PodSelector{
			Selector: chaosmesh.PodSelectorSpec{
				GenericSelectorSpec: chaosmesh.GenericSelectorSpec{
					Namespaces: []string{args.TargetNamespace},
				},
			},
			Mode: chaosmesh.AllMode,
		},
	}

}

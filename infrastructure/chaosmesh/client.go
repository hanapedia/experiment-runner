package chaosmesh

import (
	"context"
	"fmt"
	"strings"

	"github.com/hanapedia/rca-experiment-runner/infrastructure/crd/chaosmesh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type ChaosMeshClient struct {
	dynamicClient dynamic.Interface
}

func NewChaosMeshClient(dynamicClient dynamic.Interface) *ChaosMeshClient {
	return &ChaosMeshClient{
		dynamicClient: dynamicClient,
	}
}

func (client *ChaosMeshClient) ApplyNetworkDelay(networkDelay *chaosmesh.NetworkChaos) error {
	chaosMeshGVR := schema.GroupVersionResource{
		Group:    ChaosMeshGroup,
		Version:  ChaosMeshVersion,
		Resource: strings.ToLower(ChaosMeshNetworkChaosResource),
	}

	unstructuredNetworkDelay, err := runtime.DefaultUnstructuredConverter.ToUnstructured(networkDelay)
	if err != nil {
		return fmt.Errorf("could not convert networkDelay to unstructured: %w", err)
	}

	unstructuredObj := &unstructured.Unstructured{
		Object: unstructuredNetworkDelay,
	}

	// creates the network delay chaos experiment
	_, err = client.dynamicClient.Resource(chaosMeshGVR).Namespace(networkDelay.Namespace).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("could not apply the network delay chaos experiment: %w", err)
	}

	return nil
}

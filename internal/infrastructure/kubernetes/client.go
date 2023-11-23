package kubernetes

import (
	"context"
	"fmt"

	"github.com/hanapedia/experiment-runner/internal/constants"
	v1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesClient struct {
	clientset *kubernetes.Clientset
}

func NewKubernetesClient(clientset *kubernetes.Clientset) *KubernetesClient {
	return &KubernetesClient{
		clientset: clientset,
	}
}

// GetDeploymentsWithAnnotation gets list of v1.Deployment in a namespace without matching annotations
func (client *KubernetesClient) GetDeploymentsWithOutAnnotation(namespace string, annotationKey string, annotationValue string) ([]v1.Deployment, error) {
	deploymentsList, err := client.clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get deployments from namespase: %s. %s", namespace, err)
	}
	var matchingDeployments []v1.Deployment
	for _, deployment := range deploymentsList.Items {
		if deployment.Annotations[annotationKey] != annotationValue {
			matchingDeployments = append(matchingDeployments, deployment)
		}
	}
	return matchingDeployments, nil
}

// ApplyJobResource applies the batchv1.Job resource to the Cluster.
func (client *KubernetesClient) ApplyJobResource(job *batchv1.Job) error {
	_, err := client.clientset.BatchV1().Jobs(constants.RcaNamespace).Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("could not apply the job resource: %w", err)
	}
	return nil
}

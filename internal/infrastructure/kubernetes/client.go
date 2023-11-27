package kubernetes

import (
	"context"
	"fmt"

	appv1 "k8s.io/api/apps/v1"
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
func (client *KubernetesClient) GetDeploymentsWithOutAnnotation(namespace string, annotationKey string, annotationValue string) ([]appv1.Deployment, error) {
	deploymentsList, err := client.clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get deployments from namespase: %s. %s", namespace, err)
	}
	var matchingDeployments []appv1.Deployment
	for _, deployment := range deploymentsList.Items {
		if deployment.Annotations[annotationKey] != annotationValue {
			matchingDeployments = append(matchingDeployments, deployment)
		}
	}
	return matchingDeployments, nil
}

// ApplyJobResource applies the batchv1.Job resource to the Cluster.
func (client *KubernetesClient) ApplyJobResource(job *batchv1.Job, namespace string) error {
	_, err := client.clientset.BatchV1().Jobs(namespace).Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("could not apply the job resource: %w", err)
	}
	return nil
}

// CreateDeployment creates provided deployment
func (client *KubernetesClient) CreateDeployment(dep *appv1.Deployment, namespace string) error {
	_, err := client.clientset.AppsV1().Deployments(namespace).Create(context.Background(), dep, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("could not create the deployment: %w", err)
	}
	return nil
}

// DeleteDeployment deletes specified deployment
func (client *KubernetesClient) DeleteDeployment(name, namespace string) error {
	err := client.clientset.AppsV1().Deployments(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("could not delete the deployment: %w", err)
	}
	return nil
}

package cluster

import (
	"errors"
	"strings"

	"github.com/qa-kit/awesome-grid/config"
	"github.com/qa-kit/awesome-grid/deploymentconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	logger "github.com/sirupsen/logrus"
)

// Kubernetes helps to access k8s cluster methods
type Kubernetes struct {
	client dynamic.Interface
	Config *config.Config
}

// CreateClient creayes k8s client
func (k *Kubernetes) CreateClient() error {
	// creates the in-cluster config
	c, err := rest.InClusterConfig()
	if err != nil {
		return errors.New("building k8s config, " + err.Error())
	}

	client, err := dynamic.NewForConfig(c)
	if err != nil {
		return errors.New("creating dynamic k8s config, " + err.Error())
	}
	k.client = client

	return nil
}

// CreateDeployment creates deployment
func (k *Kubernetes) CreateDeployment(deploymentData deploymentconfig.DeploymentConfig) (name string, err error) {
	// Building a config
	deploymentConfig := &unstructured.Unstructured{
		Object: deploymentData,
	}
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	result, err := k.client.Resource(deploymentRes).Namespace(k.Config.Namespace).Create(deploymentConfig, metav1.CreateOptions{})
	if err != nil {
		return "", errors.New("creating k8s deployment, " + err.Error())
	}

	return result.GetName(), nil
}

// FindPodIP finds pod's ip
func (k *Kubernetes) FindPodIP(deploymentName string) (string, error) {
	// List of pods
	podRes := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	list, err := k.client.Resource(podRes).Namespace(k.Config.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return "", errors.New("listing k8s pods, " + err.Error())
	}
	for _, d := range list.Items {
		podIP, _, err := unstructured.NestedString(d.Object, "status", "podIP")
		if err != nil {
			return "", errors.New("getting k8s ip of pod, " + err.Error())
		}

		if strings.Contains(d.GetName(), deploymentName) {
			if podIP == "" {
				return "", errors.New("ip is empty for " + deploymentName)
			}
			return podIP, nil
		}
	}

	return "", errors.New("pod for deployment " + deploymentName + " not found")
}

// DeleteDeployment deletes deployment
func (k *Kubernetes) DeleteDeployment(deploymentName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	if err := k.client.Resource(deploymentRes).Namespace(k.Config.Namespace).Delete(deploymentName, deleteOptions); err != nil {
		logger.Errorf("deleting deployment %s  in k8s cluster, %s", deploymentName, err.Error())
	}

	return nil
}

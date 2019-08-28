package main

// Cluster interface provides common function for work and testing with k8s cluster
type Cluster interface {
	CreateDeployment(deploymentData DeploymentConfig) (name string, err error)
	FindPodIP(deploymentName string) (string, error)
	DeleteDeployment(deploymentName string) error
}

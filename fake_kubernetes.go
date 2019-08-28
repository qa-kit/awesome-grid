package main

// FakeKubernetes is rub for k8s client
type FakeKubernetes struct {
	CreateDeploymentResult string
	FindPodIPResult        string
}

// CreateDeployment creates deployment
func (k FakeKubernetes) CreateDeployment(deploymentData DeploymentConfig) (name string, err error) {
	return k.CreateDeploymentResult, nil
}

// FindPodIP finds pod's ip
func (k FakeKubernetes) FindPodIP(deploymentName string) (string, error) {
	return k.FindPodIPResult, nil
}

// DeleteDeployment deletes deployment
func (k FakeKubernetes) DeleteDeployment(deploymentName string) error {
	return nil
}

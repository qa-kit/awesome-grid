package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

// Creator provides creating logic of pods
type Creator struct {
	config  *Config
	cluster Cluster
	pool    *Pool
	cleaner *Cleaner
}

// Resolve creates new pod and resolve to it
func (c *Creator) Resolve(request *http.Request) (string, error) {
	// Creating config
	deploymentConfig, err := DeploymentFromTemplate(c.config.DeploymentTemplate)
	if err != nil {
		return "", errors.New("creating deployment, " + err.Error())
	}

	// Creating deployemnt
	log.Println("creating new deployment in k8s")
	name, err := c.cluster.CreateDeployment(deploymentConfig)
	if err != nil {
		return "", errors.New("creating deployment, " + err.Error())
	}
	log.Println("deployment " + name + " created")

	// Delayed removing
	go c.cleaner.DeleteDeployment(name, c.cluster, c.config.PodLifetime, c.pool)
	time.Sleep(time.Duration(c.config.WaitForCreatingTimeout) * 1000 * time.Millisecond)

	// Getting ip of new pod
	IP, err := c.cluster.FindPodIP(name)
	if err != nil {
		return "", errors.New("finding pod ip, " + err.Error())
	}

	log.Println("deployment " + name + " ip is " + IP)

	// Addind to pool
	c.pool.AddPod(name, IP)

	return "http://" + IP + ":" + c.config.PodPort, nil
}

package main

import (
	"log"
	"time"
)

// Cleaner provides delayed removing pods
type Cleaner struct {
}

// DeleteDeployment deletes deployment in cluster
func (c *Cleaner) DeleteDeployment(deploymentName string, cluster Cluster, sleep int, pool *Pool) {
	time.Sleep(time.Duration(sleep) * 1000 * time.Millisecond)
	err := cluster.DeleteDeployment(deploymentName)
	if err == nil {
		log.Println(deploymentName, "removed")
	}
	if !pool.Delete(deploymentName) {
		log.Println("deployment " + deploymentName + " not presented in local pool")
	}
}

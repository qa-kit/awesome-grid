package cleaner

import (
	"log"
	"time"

	"github.com/qa-kit/awesome-grid/cluster"
	"github.com/qa-kit/awesome-grid/pool"
)

// Cleaner provides delayed removing pods
type Cleaner struct {
}

// DeleteDeployment deletes deployment in cluster
func (c *Cleaner) DeleteDeployment(deploymentName string, cluster cluster.Cluster, sleep int, pool *pool.Pool) {
	time.Sleep(time.Duration(sleep) * 1000 * time.Millisecond)
	err := cluster.DeleteDeployment(deploymentName)
	if err == nil {
		log.Println(deploymentName, "removed")
	}
	if !pool.Delete(deploymentName) {
		log.Println("deployment " + deploymentName + " not presented in local pool")
	}
}

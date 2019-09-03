package cleaner

import (
	"time"

	"github.com/qa-kit/awesome-grid/cluster"
	"github.com/qa-kit/awesome-grid/pool"
	logger "github.com/sirupsen/logrus"
)

// Cleaner provides delayed removing pods
type Cleaner struct {
}

// DeleteDeployment deletes deployment in cluster
func (c *Cleaner) DeleteDeployment(deploymentName string, cluster cluster.Cluster, sleep int, pool *pool.Pool) {
	time.Sleep(time.Duration(sleep) * 1000 * time.Millisecond)
	err := cluster.DeleteDeployment(deploymentName)

	if err == nil {
		logger.Infof("%s removed", deploymentName)
	}

	if !pool.Delete(deploymentName) {
		logger.Infof("deployment %s not presented in local pool", deploymentName)
	}
}

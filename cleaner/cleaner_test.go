package cleaner

import (
	"testing"

	"github.com/qa-kit/awesome-grid/cluster"
	poolPkg "github.com/qa-kit/awesome-grid/pool"
)

func TestCleanerDeleteDeployment(t *testing.T) {
	var (
		deploymentName = "name1"
		host           = "127.0.0.1"
		sessionID      = "id1"
		cluster        = cluster.FakeKubernetes{}
		sleep          = 0
		pool           = poolPkg.Pool{}
		cleaner        = Cleaner{}
	)
	pool.AddPod(deploymentName, host)
	pool.AddSession(sessionID, host)
	cleaner.DeleteDeployment(deploymentName, &cluster, sleep, &pool)

	_, exists := pool.FindPodIP(deploymentName)
	if exists {
		t.Errorf("expected exists %t instead %t", exists, false)
	}
}

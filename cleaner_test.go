package main

import (
	"testing"
)

func TestCleanerDeleteDeployment(t *testing.T) {
	var (
		deploymentName = "name1"
		host           = "127.0.0.1"
		sessionID      = "id1"
		cluster        = FakeKubernetes{}
		sleep          = 0
		pool           = Pool{}
		cleaner        = Cleaner{}
	)
	pool.AddPod(deploymentName, host)
	pool.AddSession(sessionID, host)
	cleaner.DeleteDeployment(deploymentName, &cluster, sleep, &pool)

	_, exists := pool.podNameIPMap[deploymentName]
	if exists {
		t.Errorf("expected exists %t instead %t", exists, false)
	}
}

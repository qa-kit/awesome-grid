package main

import (
	"testing"
)

func TestAddPod(t *testing.T) {
	testData := map[string]string{
		"name1": "127.0.0.1",
		"name2": "127.0.0.2",
	}
	pool := Pool{}

	for name, IP := range testData {
		pool.AddPod(name, IP)
		found, exists := pool.podNameIPMap[name]
		if !exists {
			t.Errorf("expected exists %t instead %t", exists, true)
		}
		if found != IP {
			t.Errorf("expected found %s instead %s", found, IP)
		}
	}
}

func TestAddSession(t *testing.T) {
	testData := map[string]string{
		"id1": "127.0.0.1",
		"id2": "127.0.0.2",
	}
	pool := Pool{}

	for sessionID, IP := range testData {
		pool.AddSession(sessionID, IP)
		found, exists := pool.sessionIDPodIP[sessionID]
		if !exists {
			t.Errorf("expected exists %t instead %t", exists, true)
		}
		if found != IP {
			t.Errorf("expected found %s instead %s", found, IP)
		}
	}
}

func TestIP(t *testing.T) {
	testData := map[string]string{
		"id1": "127.0.0.1",
		"id2": "127.0.0.2",
	}
	pool := Pool{}

	for sessionID, IP := range testData {
		pool.AddSession(sessionID, IP)
		found, exists := pool.IP(sessionID)
		if !exists {
			t.Errorf("expected exists %t instead %t", exists, true)
		}
		if found != IP {
			t.Errorf("expected found %s instead %s", found, IP)
		}
	}
}

func TestDeleteDeployment(t *testing.T) {
	sessions := map[string]string{
		"id1": "127.0.0.1",
		"id2": "127.0.0.2",
	}
	pods := map[string]string{
		"name1": "127.0.0.1",
		"name2": "127.0.0.2",
	}
	pool := Pool{}

	for sessionID, IP := range sessions {
		pool.AddSession(sessionID, IP)
	}

	for name, IP := range pods {
		pool.AddPod(name, IP)
	}

	for name := range pods {
		result := pool.Delete(name)
		if !result {
			t.Errorf("expected exists %t instead %t", result, false)
		}
		found, exists := pool.IP(name)
		if exists {
			t.Errorf("expected exists %t instead %t", exists, false)
		}
		if found != "" {
			t.Errorf("expected found %s instead %s", found, "")
		}
	}

	result := pool.Delete("name3")
	if result {
		t.Errorf("expected exists %t instead %t", result, true)
	}
}

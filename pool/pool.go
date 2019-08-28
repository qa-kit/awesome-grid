package pool

import (
	"strings"
	"sync"
)

// Pool stores pods
type Pool struct {
	m              sync.RWMutex
	podNameIPMap   map[string]string
	sessionIDPodIP map[string]string
}

// AddPod add pod to colelction
func (p *Pool) AddPod(name string, ip string) {
	p.m.Lock()
	if p.podNameIPMap == nil {
		p.podNameIPMap = make(map[string]string)
	}
	p.podNameIPMap[name] = ip
	p.m.Unlock()
}

//FindPodIP gets pod ip by pod name
func (p *Pool) FindPodIP(name string) (string, bool) {
	p.m.RLock()
	name, exists := p.podNameIPMap[name]
	p.m.RUnlock()

	return name, exists
}

// AddSession add pod to colelction
func (p *Pool) AddSession(sessionID string, ip string) {
	p.m.Lock()
	if p.sessionIDPodIP == nil {
		p.sessionIDPodIP = make(map[string]string)
	}
	p.sessionIDPodIP[sessionID] = ip
	p.m.Unlock()
}

// IP returns pod ip by session id
func (p *Pool) IP(sessionID string) (string, bool) {
	p.m.RLock()
	ip, exists := p.sessionIDPodIP[sessionID]
	p.m.RUnlock()
	return ip, exists
}

// Delete removes data about deployment
func (p *Pool) Delete(deploymentName string) bool {
	p.m.Lock()
	defer p.m.Unlock()
	// Lookup for pod
	for k, v := range p.podNameIPMap {
		if strings.Contains(k, deploymentName) {
			delete(p.podNameIPMap, k)

			//Lookup for ip data of pod
			for z, q := range p.sessionIDPodIP {
				if q == v {
					delete(p.sessionIDPodIP, z)
					return true
				}
			}
		}
	}

	return false
}

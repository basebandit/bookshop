package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/basebandit/bookshop/server/pkg/discovery"
)

// Registry defines an in-memory service discovery registry.
type Registry struct {
	sync.RWMutex
	serviceAddrs map[string]map[string]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// NewRegistry creates a new in-memory service registry instance.
func NewRegistry() *Registry {
	return &Registry{serviceAddrs: make(map[string]map[string]*serviceInstance)}
}

// Register creates a service record in the registry.
func (r *Registry) Register(_ context.Context, instanceID, serviceName, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = make(map[string]*serviceInstance)
	}
	r.serviceAddrs[serviceName][instanceID] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(_ context.Context, instanceID, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName], instanceID)
	return nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return errors.New("service name is not registered yet")
	}
	if _, ok := r.serviceAddrs[serviceName][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}

	r.serviceAddrs[serviceName][instanceID].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(_ context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddrs[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, i := range r.serviceAddrs[serviceName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}

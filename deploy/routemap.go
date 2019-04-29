package deploy

import (
	"sync"
)

// A cache for the routes. This is almost equivalent to an ingress controller
// and is used primarily when the manager is just a Docker client. The routeMap
// maps the address (i.e. /my-project/dev) to the actual address on the machine
// (i.e. localhost:40000, kubernetes.default.svc/...)
type routeMap struct {
	m      sync.RWMutex
	routes map[string]string
}

// Set the route to the project key.
func (r *routeMap) Set(addr, route string) {
	r.m.Lock()
	defer r.m.Unlock()

	r.routes[addr] = route
}

// Get the route given the project key. If u
func (r *routeMap) Get(addr string) string {
	r.m.RLock()
	defer r.m.RUnlock()

	route, ok := r.routes[addr]
	if !ok {
		return ""
	}
	return route
}

// Remove the address from the routeMap
func (r *routeMap) Delete(addr string) {
	r.m.Lock()
	defer r.m.Unlock()

	delete(r.routes, addr)
}

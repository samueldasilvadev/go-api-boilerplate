package registry

import "fmt"

type Registry struct {
	provideFunctions map[string]func() Dependency
	deps             map[string]Dependency
}

type Dependency interface{}

func NewRegistry() *Registry {
	return &Registry{
		deps:             make(map[string]Dependency),
		provideFunctions: make(map[string]func() Dependency),
	}
}

func (r *Registry) Provide(name string, init Dependency) {
	r.deps[name] = init
}

func (r *Registry) OnDemandProvide(name string, init func() Dependency) {
	r.provideFunctions[name] = init
}

func (r *Registry) Inject(dep string) Dependency {
	depInstance, ok := r.deps[dep]
	if !ok {
		fmt.Println(dep + " Initialized")
		currentDep, found := r.provideFunctions[dep]
		if !found {
			panic("Dependency not found: " + dep)
		}
		r.deps[dep] = currentDep()
		return r.deps[dep]
	}
	return depInstance
}

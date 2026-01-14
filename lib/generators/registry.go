package generators

import "fmt"

type Registry struct {
	generators map[string]Generator
}

func NewRegistry() *Registry {
	return &Registry{
		generators: make(map[string]Generator),
	}
}

func (r *Registry) Register(a Generator) {
	r.generators[a.Name()] = a
}

func (r *Registry) Get(name string) (Generator, error) {
	a, ok := r.generators[name]
	if !ok {
		return nil, fmt.Errorf("Generator not found: %s", name)
	}
	return a, nil
}

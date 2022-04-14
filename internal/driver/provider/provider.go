package provider

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// Provider defines the behaviours of the provider
type Provider interface {
	// ContainerConfig returns the container configuration
	ContainerConfig() *container.Config
	// ContainerHostConfig returns the container host configuration
	ContainerHostConfig() *container.HostConfig
	// PortKey returns the port key
	PortKey() nat.Port
	// DataSourceName returns the data source name
	DataSourceName(port nat.PortBinding) string
}

// Factory returns the provider factory
func Factory(name string) (factory, error) {
	if factory, existed := factories[name]; existed {
		return factory, nil
	}
	return nil, fmt.Errorf("providerFactory of %q is not registered", name)
}

type factory = func(image string) Provider

var factories = make(map[string]factory)

func registerProvider(name string, providerFactory factory) {
	if providerFactory == nil {
		panic(fmt.Sprintf("providerFactory of %q is nil", name))
	}
	if _, existed := factories[name]; existed {
		panic("providerFactory of %q is already registered" + name)
	}
	factories[name] = providerFactory
}

//go:build mongo
// +build mongo

package provider

import (
	"github.com/fynxiu/dbd/internal/constant"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

func init() {
	registerProvider(constant.EngineMongo, newMongoProvider)
}

var _ Provider = (*mongoProvider)(nil)

func newMongoProvider(image string) Provider {
	return &mongoProvider{
		image: image,
	}
}

type mongoProvider struct {
	image string
}

// ContainerConfig implements provider
func (*mongoProvider) ContainerConfig() *container.Config {
	panic("unimplemented")
}

// ContainerHostConfig implements provider
func (*mongoProvider) ContainerHostConfig() *container.HostConfig {
	panic("unimplemented")
}

// DataSourceName implements provider
func (*mongoProvider) DataSourceName(port nat.PortBinding) string {
	panic("unimplemented")
}

// PortKey implements provider
func (*mongoProvider) PortKey() nat.Port {
	panic("unimplemented")
}

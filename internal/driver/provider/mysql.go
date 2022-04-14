package provider

import (
	"fmt"

	"github.com/fynxiu/dbd/internal/constant"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

func init() {
	registerProvider(constant.EngineMysql, newMysqlProvider)
}

const (
	mysqlProtoclPort  = "3306/tcp"
	mysqlRootPassword = "20140119"
	mysqlDatabaseName = constant.DatabaseName
)

var _ Provider = (*mysqlProvider)(nil)

type mysqlProvider struct {
	image string
}

func newMysqlProvider(image string) Provider {
	return &mysqlProvider{
		image: image,
	}
}

// ContainerConfig implements provider
func (p *mysqlProvider) ContainerConfig() *container.Config {
	return &container.Config{
		Image: p.image,
		ExposedPorts: nat.PortSet{
			mysqlProtoclPort: {},
		},
		Env: []string{
			"MYSQL_ROOT_PASSWORD=" + mysqlRootPassword,
			"MYSQL_DATABASE=" + mysqlDatabaseName,
		},
	}
}

// ContainerHostConfig implements provider
func (*mysqlProvider) ContainerHostConfig() *container.HostConfig {
	return &container.HostConfig{
		PortBindings: nat.PortMap{
			mysqlProtoclPort: []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				},
			},
		},
	}
}

// PortKey implements provider
func (*mysqlProvider) PortKey() nat.Port {
	return mysqlProtoclPort
}

// DataSourceName implements provider
func (*mysqlProvider) DataSourceName(port nat.PortBinding) string {
	return fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8mb4&interpolateParams=true", mysqlRootPassword, port.HostIP, port.HostPort, mysqlDatabaseName)
}

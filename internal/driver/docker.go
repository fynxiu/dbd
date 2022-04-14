package driver

import (
	"context"

	"github.com/fynxiu/dbd/internal/constant"
	"github.com/fynxiu/dbd/internal/driver/provider"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
)

var _ Driver = (*dockerDriver)(nil)

const (
	containerName = "fyn_sql_deduction_container"
	filterKeyName = "name"
)

// NewDockerDriver defines a docker implementation of Driver
func NewDockerDriver(provider provider.Provider, image string) (Driver, error) {
	client, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv)
	if err != nil {
		return nil, err
	}
	return &dockerDriver{
		provider: provider,
		client:   client,
		image:    image,
	}, nil

}

type dockerDriver struct {
	provider    provider.Provider
	client      *dockerclient.Client
	image       string
	portBinding nat.PortBinding
}

// Reuse implements Driver
func (d *dockerDriver) Reuse(ctx context.Context) error {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   filterKeyName,
			Value: containerName,
		}),
	})
	if err != nil {
		return err
	}
	if len(containers) == 0 {
		return constant.ErrContainerNotFound
	}

	return d.inspectPortBinding(ctx, containers[0].ID)
}

// Launch implements Driver
func (d *dockerDriver) Launch(ctx context.Context) error {
	if err := d.ensureCleanEnvironment(ctx); err != nil {
		return nil
	}
	resp, err := d.client.ContainerCreate(ctx, d.provider.ContainerConfig(),
		d.provider.ContainerHostConfig(),
		nil,
		nil, containerName)
	if err != nil {
		return err
	}

	cid := resp.ID
	if err := d.client.ContainerStart(ctx, cid, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return d.inspectPortBinding(ctx, cid)
}

func (d *dockerDriver) inspectPortBinding(ctx context.Context, cid string) error {
	inspect, err := d.client.ContainerInspect(ctx, cid)
	if err != nil {
		return err
	}
	d.portBinding = inspect.NetworkSettings.Ports[d.provider.PortKey()][0]
	return nil
}

// DataSourceName implements Driver
func (d *dockerDriver) DataSourceName() string {
	return d.provider.DataSourceName(d.portBinding)
}

// Dispose implements Driver
func (d *dockerDriver) Dispose(ctx context.Context) error {
	return d.ensureCleanEnvironment(ctx)
}

// ensureCleanEnvironment ensures the environment is clean
func (d *dockerDriver) ensureCleanEnvironment(ctx context.Context) error {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   filterKeyName,
			Value: containerName,
		}),
	})
	if err != nil {
		return err
	}
	for i := range containers {
		if err := d.client.ContainerRemove(ctx, containers[i].ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			return err
		}
	}
	return nil
}

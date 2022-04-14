package driver

import (
	"context"

	"github.com/fynxiu/dbd/internal/constant"
	"github.com/fynxiu/dbd/internal/driver/provider"
)

type (
	// Driver defines the behaviours of the driver
	Driver interface {
		// Reuse reuses the driver
		Reuse(ctx context.Context) error
		// Launch launches the driver
		Launch(ctx context.Context) error
		DataSourceName() string
		Dispose(context.Context) error
	}
)

// NewDefaultDriver default factory method for Driver
func NewDefaultDriver(engine, image string) (Driver, error) {
	providerFactory, err := provider.Factory(engine)
	if err != nil {
		return nil, constant.ErrEngineNotSupported
	}

	return NewDockerDriver(providerFactory(image), image)
}

package constant

import "errors"

// ErrEngineNotSupported is returned when the engine is not supported.
var ErrEngineNotSupported = errors.New("engine not supported")

// ErrContainerNotFound is returned when the container is not found.
var ErrContainerNotFound = errors.New("container not found")

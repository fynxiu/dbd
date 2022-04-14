package executor

import (
	"fmt"
)

type executorFactory = func(DataSourceName string) Executor

// Factories is the map of executor factories.
var Factories = make(map[string]executorFactory)

// Executor is the interface that wraps the ExecuteScript and Schema methods.
type Executor interface {
	// ExecuteScript execute the script.
	ExecuteScript(script string) error
	// Schema returns the schema of the database.
	Schema() (string, error)
	// Dispose closes the database connection.
	Dispose() error
}

func registerExecutor(name string, executorFactory executorFactory) {
	if executorFactory == nil {
		panic(fmt.Sprintf("executorFactory of %q is nil", name))
	}
	if _, existed := Factories[name]; existed {
		panic("executorFactory of %q is already registered" + name)
	}
	Factories[name] = executorFactory
}

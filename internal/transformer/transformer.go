package transformer

import (
	"fmt"
	"io/ioutil"

	"github.com/fynxiu/dbd/internal/constant"
	"github.com/fynxiu/dbd/internal/transformer/executor"
)

type (
	// Transformer defines the interface for transforming
	Transformer interface {
		// Transform transforms the given files into a schema
		Transform(files []string) (string, error)
	}
)

// NewTransformer creates a new Transformer
func NewTransformer(engine string, dataSourceName string) (Transformer, error) {
	executorFactory, ok := executor.Factories[engine]
	if !ok {
		return nil, constant.ErrEngineNotSupported
	}

	return &defaultTransformer{
		executorFactory(dataSourceName),
	}, nil
}

var _ Transformer = (*defaultTransformer)(nil)

type defaultTransformer struct {
	executor executor.Executor
}

// Transform implements Transformer
func (t *defaultTransformer) Transform(files []string) (string, error) {
	defer t.executor.Dispose()

	for _, x := range files {
		script, err := ioutil.ReadFile(x)
		if err != nil {
			return "", err
		}

		if err := t.executor.ExecuteScript(string(script)); err != nil {
			return "", fmt.Errorf("Transform failed, %v", err)
		}
	}

	return t.executor.Schema()
}

package config

import (
	"fmt"

	"github.com/fynxiu/dbd/internal/constant"
)

// EngineConfig returns the configuration for the given engine type.
type EngineConfig struct {
	Image          string
	OutputFilename string
	Comment        func(string) string
}

var configMap = map[string]EngineConfig{
	constant.EngineMysql: {
		"mysql:8.0.24",
		"output.sql",
		func(s string) string { return fmt.Sprintf("-- %s", s) },
	},
	constant.EngineMongo: {
		"mongo:latest",
		"output.js",
		func(s string) string { return fmt.Sprintf("// %s", s) },
	},
}

// GetEngineConfig returns the configuration for the given engine type.
// panic if the engine type is not supported.
func GetEngineConfig(engineType string) *EngineConfig {
	if config, ok := configMap[engineType]; ok {
		return &config
	}
	panic(fmt.Sprintf("unknown engine type: %s", engineType))
}

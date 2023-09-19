package config

import (
	"strings"
)

// BackendConfig represents the configuration for backends.
type BackendConfig map[string]string

func LoadBackends(args []string) BackendConfig {
	argsMap := make(map[string]string)

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)

		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			argsMap[key] = value
		}
	}

	return argsMap
}

package config

import (
	"fmt"
	"os"
)

// Env checks if env var or its default values exists and returns as string value.
func Env(envVar, defaultEnvVar string) (string, error) {
	v := os.Getenv(envVar)

	if v == "" && defaultEnvVar != "" {
		return defaultEnvVar, nil
	}

	if v == "" && defaultEnvVar == "" {
		return "", fmt.Errorf("missing env var / default env var: %v", envVar)
	}

	return v, nil
}

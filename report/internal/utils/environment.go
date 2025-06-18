// Package envar gets data from environment vars
package utils

import "os"

// GetEnvVar the key from environment, if not set or empty returns def.
func GetEnvVar(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}

package envar

import "os"

// Get the key from environment, if not set or empty returns def.
func Get(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}

package envx

import "os"

// Get works like os.Getenv, but adds the ability ro povide a default
func Get(key string, def string) (value string) {
	value = def
	if v := os.Getenv(key); v != "" {
		value = v
	}
	return
}

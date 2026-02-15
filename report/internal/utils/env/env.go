package env

import (
	"opg-reports/report/internal/utils/marshal"
	"os"
	"strings"
)

// Get works like os.Getenv, but adds the ability ro povide a default
func Get(key string, def string) (value string) {
	value = def
	if v := os.Getenv(key); v != "" {
		value = v
	}
	return
}

// OverwriteStruct will overwrite values in the struct with those found
// from os.Getenv.
//
// Uses the uppercase version of the key name (so `id` => `ID`) and and
// hyphens become underscores.
func OverwriteStruct[T any](data T) (err error) {
	var json = map[string]interface{}{}
	// convert data to a map for checking
	err = marshal.Convert(data, &json)
	if err != nil {
		return
	}
	// check for each uppercase version of the key name
	for key, _ := range json {
		osKey := strings.ReplaceAll(strings.ToUpper(key), "-", "_")
		if v := os.Getenv(osKey); v != "" {
			json[key] = v
		}
	}
	// convert back
	err = marshal.Convert(json, &data)
	return
}

package structs

import (
	"encoding/json"
	"log/slog"
)

// ToMap takes a struct, marshals that to json and then unmarshals
// that into a map[string]interface.
// Used to get values out of structs etc
func ToMap[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	} else {
		slog.Error("[structs] failed to covnert to a map", slog.String("err", err.Error()))
	}
	return
}

func Unmarshal[T any](content []byte, destination T) (err error) {
	err = json.Unmarshal(content, destination)
	return
}

package convert

import (
	"encoding/json"
	"log/slog"
)

func Map[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	} else {
		slog.Error("map failed", slog.String("err", err.Error()))
	}
	return
}

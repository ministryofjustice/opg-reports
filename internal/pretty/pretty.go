package pretty

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

func Print[T any](item T) {
	var err error
	var bytes []byte
	if bytes, err = json.MarshalIndent(item, "", "  "); err == nil {
		fmt.Println(string(bytes))
	}
	if err != nil {
		slog.Error("[pretty] print", slog.String("err", err.Error()))
	}
}

package pretty

import (
	"encoding/json"
	"fmt"
)

func Print[T any](item T) {
	var err error
	var bytes []byte
	if bytes, err = json.MarshalIndent(item, "", "  "); err == nil {
		fmt.Println(string(bytes))
	}
}

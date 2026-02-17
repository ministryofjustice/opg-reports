package dump

import (
	"encoding/json"
	"fmt"
)

func Any(item any) string {
	var str = ""
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err == nil {
		str = string(bytes)
	}
	return fmt.Sprintf("%+v\n", str)
}

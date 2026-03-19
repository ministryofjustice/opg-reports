package fmtx

import (
	"encoding/json"
	"fmt"
)

func Dump(item any) {
	var str = ""
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err == nil {
		str = string(bytes)
	}
	fmt.Println(fmt.Sprintf("%+v\n", str))
}

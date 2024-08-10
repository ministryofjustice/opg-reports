package convert

import (
	"encoding/json"
)

func ToJson[T any](item T) (content []byte, err error) {
	return json.MarshalIndent(item, "", "  ")
}

func ListToJson[T any](items []T) (content []byte, err error) {
	return json.MarshalIndent(items, "", "  ")
}

func ToMap[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	}
	return
}

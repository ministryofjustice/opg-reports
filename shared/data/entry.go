package data

import "encoding/json"

// IEntry used to represent an item within the data store
// and result from a report
type IEntry interface {
	Valid() bool
	UID() string
}

// ToMap uses json marshaling to convert from the struct to map.
// Does require struct to be tagged correctly to do this neatly
func ToMap[T IEntry](item T) (m map[string]string, err error) {
	jBytes, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(jBytes, &m)
	}
	return
}

// FromMap uses json marshaling to convert from a map back to a struct.
// Requires the struct to be tagged correctly to match fields etc
func FromMap[T IEntry](m map[string]string) (item T, err error) {
	jBytes, err := json.Marshal(m)
	if err == nil {
		json.Unmarshal(jBytes, &item)
	}
	return
}

// FromJson will convert a []byte (normally result of file stream or
// json marshalling) back to the struct.
// Does required corect tagging on the struct.
func FromJson[T IEntry](content []byte) (item T, err error) {
	err = json.Unmarshal(content, &item)
	return
}

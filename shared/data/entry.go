package data

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

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
	slog.Debug("[data/entry] ToMap", slog.String("UID", item.UID()), slog.String("err", fmt.Sprintf("%v", err)))
	return
}

// ToJson converts item T to a []byte via json marshaling
func ToJson[T IEntry](item T) (content []byte, err error) {
	content, err = json.MarshalIndent(item, "", "  ")
	slog.Debug("[data/entry] ToJson", slog.String("UID", item.UID()), slog.String("err", fmt.Sprintf("%v", err)))
	return
}

// ToJsonList converts a series of T into a []byte string via marshalling
func ToJsonList[T IEntry](items []T) (content []byte, err error) {
	content, err = json.MarshalIndent(items, "", "  ")
	slog.Debug("[data/entry] ToJsonList", slog.String("err", fmt.Sprintf("%v", err)))
	return
}

// FromMap uses json marshaling to convert from a map back to a struct.
// Requires the struct to be tagged correctly to match fields etc
func FromMap[T IEntry](m map[string]string) (item T, err error) {
	jBytes, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(jBytes, &item)
	}
	slog.Debug("[data/entry] FromMap", slog.String("UID", item.UID()), slog.String("err", fmt.Sprintf("%v", err)))
	return
}

// FromJson will convert a []byte (normally result of file stream or
// json marshalling) back to the struct.
// Does required corect tagging on the struct.
func FromJson[T IEntry](content []byte) (item T, err error) {
	err = json.Unmarshal(content, &item)
	slog.Debug("[data/entry] FromJson", slog.String("UID", item.UID()), slog.String("err", fmt.Sprintf("%v", err)))
	return
}

// FromJsonList returns a slice of T ([]T) rather than a single T
func FromJsonList[T IEntry](content []byte) (items []T, err error) {
	err = json.Unmarshal(content, &items)
	slog.Debug("[data/entry] FromJsonList", slog.String("err", fmt.Sprintf("%v", err)))
	return
}

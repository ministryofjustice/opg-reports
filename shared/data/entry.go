package data

import (
	"encoding/json"
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
	return
}

// ToJson converts item T to a []byte via json marshaling
func ToJson[T IEntry](item T) (content []byte, err error) {
	return json.Marshal(item)
}

// ToJsonList converts a series of T into a []byte string via marshalling
func ToJsonList[T IEntry](items []T) (content []byte, err error) {
	return json.Marshal(items)
}

// ToInterface converts an instance of T into a map interface - generally
// used to test conversion back and forth in result handling
func ToInterface[T IEntry](item T) (iItem map[string]interface{}, err error) {
	iItem = map[string]interface{}{}

	if mapped, err := ToMap(item); err == nil {
		for key, value := range mapped {
			iItem[key] = value
		}
	}
	return
}

// ToInterfaces converts a list of T into a slice of interfaces - generally
// used to test conversion back and forth in result handling
func ToInterfaces[T IEntry](items []T) (iItems []interface{}, err error) {
	iItems = []interface{}{}
	for _, item := range items {
		if iItem, e := ToInterface[T](item); err == nil {
			iItems = append(iItems, iItem)
		} else {
			err = e
		}
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

// FromJsonList returns a slice of T ([]T) rather than a single T
func FromJsonList[T IEntry](content []byte) (items []T, err error) {
	err = json.Unmarshal(content, &items)
	return
}

// FromInterface converts an interfaces (likely from an apiresponse.Results() )
// into a T
// Presumes that the intem passed can be swapped to map[string]interface{}
// and should be converted to a map[string]string before T
func FromInterface[T IEntry](inter interface{}) (item T, err error) {
	mappedI := inter.(map[string]interface{})
	mapped := map[string]string{}
	for key, val := range mappedI {
		mapped[key] = val.(string)
	}
	if tItem, e := FromMap[T](mapped); e == nil {
		item = tItem
	} else {
		err = e
	}
	return

}

// FromInterfaces converts a series of interfaces (likely from an apiresponse.Results() )
// into a slice of T
// Presumes that each element within the interItems is a map[string]interface{}
// and should be converted to a map[string]string
func FromInterfaces[T IEntry](interItems []interface{}) (items []T, err error) {
	items = []T{}

	for _, interItem := range interItems {
		tItem, e := FromInterface[T](interItem)
		if e != nil {
			return nil, e
		}
		items = append(items, tItem)

	}
	return
}

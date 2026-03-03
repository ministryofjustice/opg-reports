package rest

import (
	"fmt"
	"net/http"
	"strings"
)

type ParamType string

const (
	QUERY ParamType = "query"
	PATH  ParamType = "path"
)

// Param
type Param struct {
	Key    string
	Type   ParamType
	Value  string
	Locked bool // stops value being replaced by request
}

// ID is used ot meet the overwrite.keyed interface requirement
func (self *Param) ID() string {
	return self.Key
}
func (self *Param) PathKey() string {
	return fmt.Sprintf("{%s}", self.Key)
}

// GetValue will give the value for this param to use in api query by checking for it
// in the current front end server request - allowing pass through
func (self *Param) GetValue(current *http.Request) (v string) {
	var values = current.URL.Query()
	v = self.Value

	// if locked, then cant be overwrittern by request
	if self.Locked {
		return self.Value
	}

	for key, set := range values {
		if key == self.Key {
			v = strings.TrimSuffix(strings.Join(set, ","), ",")
		}
	}

	return
}

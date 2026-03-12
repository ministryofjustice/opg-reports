package convert

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Between takes original struct of T and by marshaling and then unmarshaling applied its
// content to destination R
func Between[R any](source any, destination R) (err error) {
	var bytes []byte

	switch any(source).(type) {
	// convert response in a particualr way
	case *http.Response:
		var resp = source.(*http.Response)
		defer resp.Body.Close()
		bytes, err = io.ReadAll(resp.Body)
	default:
		bytes, err = json.MarshalIndent(source, "", "  ")
		if err != nil {
			return
		}
	}
	err = json.Unmarshal(bytes, destination)
	return
}

// String provides helper to wrap interface to string conversions
//
// conversion is done as follows:
//
//   - string: as is
//   - int, int32, int64: uses `strvconv.Itoa`
//   - float, float32, float64: uses `fmt.Sprintf("%g")`
//   - *http.Response: closes buffers and converts bytes to string
func String(item any) (s string) {
	s = ""

	switch any(item).(type) {
	case string:
		s = item.(string)
	case int, int32, int64:
		s = strconv.Itoa(item.(int))
	case float32, float64:
		s = fmt.Sprintf("%g", item)
	// case time.Time:
	// 	s = item.(time.Time).Format(times.YM)
	case *http.Response:
		r := item.(*http.Response)
		defer r.Body.Close()
		if content, err := io.ReadAll(r.Body); err == nil {
			s = string(content)
		}
	}
	return
}

package convert

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"opg-reports/report/packages/times"
	"strconv"
	"time"
)

// Between takes original struct of T and by marshaling and then unmarshaling applied its
// content to destination R
func Between[T any, R any](source T, destination R) (err error) {
	var bytes []byte
	if bytes, err = json.MarshalIndent(source, "", "  "); err == nil {
		err = json.Unmarshal(bytes, destination)
	}
	return
}

// Response is a specific helper to handle dealing with
// http responses and clears buffers etc
func Response[R any](response *http.Response, destination R) (err error) {
	var content []byte

	defer response.Body.Close()
	content, err = io.ReadAll(response.Body)

	if err != nil {
		fmt.Println("close error...")
		return
	}
	err = json.Unmarshal(content, destination)
	fmt.Printf("-->%v\n", err)
	return
}

// String provides helper to wrap interface to string conversions
//
// conversion is done as follows:
//
//   - string: as is
//   - int, int32, int64: uses `strvconv.Itoa`
//   - float, float32, float64: uses `fmt.Sprintf("%g")`
//   - time.Time: uses year-month (YYYY-MM) formatted string
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
	case time.Time:
		s = item.(time.Time).Format(times.YM)
	case *http.Response:
		r := item.(*http.Response)
		defer r.Body.Close()
		if content, err := io.ReadAll(r.Body); err == nil {
			s = string(content)
		}
	}
	return
}

// Ptr converts an item to the pointer version of itself.
func Ptr[T any](item T) *T {
	var ptr = &item
	return ptr
}

package cnv

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var ErrUnsupportedType = errors.New("unsupported type")

// Convert takes original struct of T and by marshaling and then unmarshaling applied its
// content to destination R
func Convert[T any, R any](source T, destination R) (err error) {
	var bytes []byte
	if bytes, err = json.MarshalIndent(source, "", "  "); err == nil {
		err = json.Unmarshal(bytes, destination)
	}
	return
}

func Capitalize(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = cases.Title(language.English).String(word)
	}
	return strings.Join(words, " ")

}

// ToFloat converts strings, ints & floats to a float64
func ToFloat(src interface{}) (dest float64, err error) {

	switch src.(any).(type) {
	case string:
		f, e := strconv.ParseFloat(src.(string), 64)
		if e != nil {
			err = e
		} else {
			dest = f
		}
	case float64:
		dest = src.(float64)
	case float32:
		dest = float64(src.(float32))
	case int:
		dest = float64(src.(int))
	case int32:
		dest = float64(src.(int32))
	case int64:
		dest = float64(src.(int64))
	default:
		err = ErrUnsupportedType
	}
	return
}

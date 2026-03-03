package cnv

import (
	"encoding/json"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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

package utils

import (
	"html/template"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Title(s string) string {
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "_", " ")
	c := cases.Title(language.English)
	s = c.String(s)
	return s
}

// Currency displays s as a float with 2 decimals and appends the currency symbol to the start
// Used to display financial values cleanly
func Currency(s string, symbol string) string {
	if s == "" {
		return ""
	}
	return symbol + FloatString(s, "%.2f")
}

// FloatString treats s as a float (if its a string, it converts, other types are ignored) and
// returns sprintf version using the layout passed
// defaults to "0.00"
func FloatString(s string, layout string) string {
	p := message.NewPrinter(language.English)
	f, _ := strconv.ParseFloat(s, 64)
	return p.Sprintf(layout, f)

}

// ValueFromMap uses the `name` as the key in the map `data` and returns
// its value
// if name is not a key on the map, it returns empty string
func ValueFromMap(name string, data map[string]string) (value string) {
	value = ""
	if v, ok := data[name]; ok {
		value = v
	}
	return
}

func TemplateFunctions() (funcs template.FuncMap) {
	funcs = map[string]interface{}{
		// simple strings
		"Title":   Title,
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
		// strinng -> numbers
		"Currency": Currency,
		// accessing maps
		"ValueFromMap": ValueFromMap,
		//
	}

	return
}

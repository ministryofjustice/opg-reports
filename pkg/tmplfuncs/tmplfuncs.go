// Package tmplfuncs contains series of funcs to use within
// front end templates
//
// Exposed as a map (`All`)
package tmplfuncs

import (
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-reports/pkg/navigation"
	"github.com/ministryofjustice/opg-reports/pkg/nums"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Increment returns a value 1 higher than the value
// passed
func Increment(i interface{}) (result interface{}) {
	result = i
	switch i.(type) {
	case float64:
		result = i.(float64) + 1
	case int:
		result = i.(int) + 1
	}

	return
}

// Title generates a title case string
func Title(s string) string {
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	c := cases.Title(language.English)
	s = c.String(s)
	return s
}

// ValueFromMap uses the `name` as the key in the map `data` and returns
// its value
// if name is not a key on the map, it returns empty string
func ValueFromMap(name string, data map[string]interface{}) (value interface{}) {
	value = ""
	if v, ok := data[name]; ok {
		value = v
	}
	return
}

func Currency(s interface{}, symbol string) string {
	return symbol + FloatString(s, "%.2f")
}

func FloatString(s interface{}, layout string) string {
	p := message.NewPrinter(language.English)
	switch s.(type) {
	case string:
		f, _ := strconv.ParseFloat(s.(string), 10)
		return p.Sprintf(layout, f)
	case float64:
		return p.Sprintf(layout, s.(float64))
	}
	return "0.00"
}

func Percentage(s interface{}) string {
	return FloatString(s, "%.4f") + " %"
}

func Matches(a *navigation.Navigation, b *navigation.Navigation) (m bool) {
	m = false
	if b != nil {
		m = a.Uri == b.Uri
	}
	return
}

var All map[string]interface{} = map[string]interface{}{
	// access
	"valueFromMap": ValueFromMap,
	"matches":      Matches,
	// numeric
	"add":       nums.Add,
	"increment": Increment,
	// formatting
	"title":      Title,
	"currency":   Currency,
	"percentage": Percentage,
}

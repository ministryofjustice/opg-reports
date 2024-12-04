// Package tmplfuncs contains series of funcs to use within
// front end templates
//
// Exposed as a map (`All`)
package tmplfuncs

import (
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-reports/info"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/internal/navigation"
	"github.com/ministryofjustice/opg-reports/internal/nums"
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
	// strip elements of the field name
	s = strings.ReplaceAll(s, "aws_account_number", "account")
	s = strings.ReplaceAll(s, "aws_account_", "")
	s = strings.ReplaceAll(s, "_name", "")

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

// Currency displays s as a float with 2 decimals and appends the currency symbol to the start
// Used to display financial values cleanly
func Currency(s interface{}, symbol string) string {
	return symbol + FloatString(s, "%.2f")
}

// FloatString treats s as a float (if its a string, it converts, other types are ignored) and
// returns sprintf version using the layout passed
// defaults to "0.00"
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

// Percentage presumes the string is a float and
// returns a string with 4 decimals and a trailing %
// marker
func Percentage(s interface{}) string {
	return FloatString(s, "%.4f") + " %"
}

// Matches is used by top navigation to see if it is a match for the active root
// section of the website - which when true should then highlight the
// nav item
// Compares the .Uri properties
func Matches(a *navigation.Navigation, b *navigation.Navigation) (m bool) {
	m = false
	if b != nil {
		m = a.Uri == b.Uri
	}
	return
}

// Day converts the string into a date and then back to a string
// formatted as yyyy-mm-dd
func Day(s string) (date string) {
	date = s
	if v, err := dateutils.Time(s); err == nil {
		date = v.Format(dateformats.YMD)
	}
	return
}

// DayBefore converts the string into a date, removes 1 day
// and then converts back to a string formatted as yyyy-mm-dd
func DayBefore(s string) (date string) {
	date = s
	if v, err := dateutils.Time(s); err == nil {
		date = v.AddDate(0, 0, -1).Format(dateformats.YMD)
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
	"day":        Day,
	"dayBefore":  DayBefore,
	// info
	"awsBillingDay": func() int { return info.AwsBillingDay },
}

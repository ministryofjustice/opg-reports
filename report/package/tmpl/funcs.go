package tmpl

import (
	"fmt"
	"opg-reports/report/package/tabulate"
	"opg-reports/report/package/times"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	billingStable   string = "stable"
	billingUnstable string = "unstable"
)

func LengthClass(i int) string {
	if i >= 35 {
		return "xl"
	} else if i >= 25 {
		return "l"
	} else if i >= 15 {
		return "m"
	}
	return "s"
}

// StripOrgansiation is used to remove github org from repository names
// to make them shorter to display
func StripOrgansiation(s string) string {
	if pos := strings.Index(s, "/"); pos >= 0 {
		s = s[pos+1:]
	}
	return s
}

func Title(s string) string {
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "_", " ")
	c := cases.Title(language.English)
	s = c.String(s)
	return s
}

// Currency displays s as a float with 2 decimals and appends the currency symbol to the start
// Used to display financial values cleanly
func Currency(s interface{}, symbol string) (c string) {
	var pntr = message.NewPrinter(language.English) // will add 1000s seperator

	switch s.(any).(type) {
	case string:
		if p, e := strconv.ParseFloat(s.(string), 64); e == nil {
			c = pntr.Sprintf("%s%.2f", symbol, p)
		}
	case float64:
		c = pntr.Sprintf("%s%.2f", symbol, s.(float64))
	}

	return
}

func Percentage(s interface{}, places int) string {
	var layout = "%." + strconv.Itoa(places) + "f"

	switch s.(any).(type) {
	case string:
		if f, e := strconv.ParseFloat(s.(string), 64); e == nil {
			return fmt.Sprintf(layout, f) + `%`
		}
	case float64:
		return fmt.Sprintf(layout, s) + `%`
	}
	return ""
}

func Headers(name string, data map[tabulate.ColType][]string) (value []string) {
	value = []string{}
	if v, ok := data[tabulate.ColType(name)]; ok {
		value = v
	}
	return
}

// ValueFromMap uses the `name` as the key in the map `data` and returns
// its value as interface
// if name is not a key on the map, it returns empty string
func ValueFromMap(name string, data map[string]interface{}) (value interface{}) {
	value = ""
	if v, ok := data[name]; ok {
		value = v
	}
	return
}
func PopRow(rows []map[string]interface{}) (foot map[string]interface{}) {
	foot = rows[len(rows)-1]
	return
}
func TrimLastRow(src []map[string]interface{}) (rows []map[string]interface{}) {
	rows = src[0 : len(src)-1]
	return
}

// BillingStabilityClass decides in the month has passed its billing date period
//
// Will return either "billing-stable" or "billing-unstable". Returns empty on error
func BillingStabilityClass(billingDay int, yymm string) (className string) {

	month, err := times.AsTime(yymm)
	if err != nil {
		return
	}

	billingMonth := times.LastBillingMonth(billingDay)
	if month.After(billingMonth) {
		className = billingUnstable
	} else {
		className = billingStable
	}

	return
}

func BillingStabilitySuffix(className string) string {
	if className == billingUnstable {
		return "*"
	}
	return ""
}

func TemplateFunctions() (funcs template.FuncMap) {
	funcs = map[string]interface{}{
		// simple strings
		"Title":             Title,
		"ToLower":           strings.ToLower,
		"ToUpper":           strings.ToUpper,
		"StripOrgansiation": StripOrgansiation,
		"LengthClass":       LengthClass,
		// string -> numbers
		"Currency":   Currency,
		"Percentage": Percentage,
		// accessing maps
		"ValueFromMap": ValueFromMap,
		"Headers":      Headers,
		"PopRow":       PopRow,
		"TrimLastRow":  TrimLastRow,
		//
		"BillingStabilityClass":  BillingStabilityClass,
		"BillingStabilitySuffix": BillingStabilitySuffix,
	}

	return
}

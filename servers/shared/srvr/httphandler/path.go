package httphandler

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

// replacerFunc signature for functions that will replace dynamic segments of the url
//
//   - url: the original url
//   - matched: the segment that has been matched (so {month} | {month:-1})
//   - modifier: the modification value (-1 in the case of {month:-1})
//   - when: pointer to time used for subs
//
// The returned string should be the original url with the match replaced with the parsed
// value
type replacerFunc func(url string, matched string, modifier string, when *time.Time) string

// Month replaces {month} | {month:<modifier>} patterns with a real value
//
// Implements replaceFunc
func Month(url string, matched string, modifier string, when *time.Time) (newUrl string) {
	newUrl = url
	date := dates.ResetMonth(*when)

	// If the modifier is a number we can use, then parse and change the month count by that
	// value
	if months, err := strconv.Atoi(modifier); err == nil && modifier != "" {
		date = date.AddDate(0, months, 0)
	}
	newUrl = strings.ReplaceAll(url, matched, date.Format(dates.FormatYM))
	return
}

// Day replaces {day} | {day:<modifier>} patterns with a real value
//
// Implements replaceFunc
func Day(url string, matched string, modifier string, when *time.Time) (newUrl string) {
	newUrl = url
	date := dates.ResetDay(*when)
	// If the modifier is a number we can use, then parse and change the day count by that
	// value
	if days, err := strconv.Atoi(modifier); err == nil && modifier != "" {
		date = date.AddDate(0, 0, days)
	}
	newUrl = strings.ReplaceAll(url, matched, date.Format(dates.FormatYMD))
	return
}

// BillingMonth replaces {billingMonth} | {billingMonth:<modifier>} patterns with a real value
//
// Implements replaceFunc
func BillingMonth(url string, matched string, modifier string, when *time.Time) (newUrl string) {
	newUrl = url
	date := dates.BillingEndDate(*when, consts.BILLING_DATE)
	// If the modifier is a number we can use, then parse and change the month count by that
	// value
	if months, err := strconv.Atoi(modifier); err == nil && modifier != "" {
		date = date.AddDate(0, months, 0)
	}
	newUrl = strings.ReplaceAll(url, matched, date.Format(dates.FormatYM))
	return
}

// BillingDay replaces {billingDay} | {billingDay:<modifier>} patterns with a real value
//
// Implements replaceFunc
func BillingDay(url string, matched string, modifier string, when *time.Time) (newUrl string) {
	newUrl = url
	date := dates.BillingEndDate(*when, consts.BILLING_DATE)
	// If the modifier is a number we can use, then parse and change the month count by that
	// value
	if days, err := strconv.Atoi(modifier); err == nil && modifier != "" {
		date = date.AddDate(0, 0, days)
	}
	newUrl = strings.ReplaceAll(url, matched, date.Format(dates.FormatYMD))
	return
}

// FuncNameAndModifier takes a matched string ({<func>:<modifier>}) and returns
// the values for <func> and <modifier>
func FuncNameAndModifier(match string) (funcName string, modifier string) {
	modifier = ""
	funcName = strings.ReplaceAll(strings.ReplaceAll(match, "}", ""), "{", "")
	// if the funcName contains : then remove that and use as modified
	if strings.Contains(funcName, ":") {
		split := strings.Split(funcName, ":")
		funcName = split[0]
		modifier = split[1]
	}
	return
}

var functionMap = map[string]replacerFunc{
	"month":        Month,
	"billingMonth": BillingMonth,
	"day":          Day,
	"billingDay":   BillingDay,
}

// Path takes the original string and converts it to a parsed version
// with the known {<func>:<modifier>} replaced with real versions
//
// Uses a regex to find all patterns that match `{.*}`, parses
// each of those looking for a matching function name from functionMap
// and calls the replacer function to swap with real values
func Path(url string) (parsed string) {
	var (
		re       = regexp.MustCompile(`(?mi){(.*?)}`)
		now      = time.Now().UTC()
		original = url
	)
	parsed = url
	//
	for _, match := range re.FindAllString(original, -1) {
		funcName, modifier := FuncNameAndModifier(match)

		if replacer, ok := functionMap[funcName]; ok {
			parsed = replacer(parsed, match, modifier, &now)
		}
		slog.Debug("[front] converted url path to rela values",
			slog.String("original", original),
			slog.String("match", match),
			slog.String("funcName", funcName),
			slog.String("modifier", modifier))
	}

	return
}

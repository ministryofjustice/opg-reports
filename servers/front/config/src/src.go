package src

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

type ApiUrl string

func (u *ApiUrl) Parsed() (s string) {
	var re = regexp.MustCompile(`(?mi){(.*?)}`)
	var now = time.Now().UTC()
	var original = string(*u)
	s = original

	for _, match := range re.FindAllString(original, -1) {
		idx, frag := indexAndFragment(match)

		if subFunc, ok := functions[idx]; ok {
			s = subFunc(match, frag, s, &now)
		}

		slog.Debug("[src] api url processing",
			slog.String("match", match),
			slog.String("idx", idx),
			slog.String("frag", frag),
			slog.String("original", original))
	}

	return
}

func indexAndFragment(match string) (index string, fragment string) {
	fragment = ""
	index = strings.ReplaceAll(strings.ReplaceAll(match, "}", ""), "{", "")
	if strings.Contains(index, ":") {
		sp := strings.Split(index, ":")
		index = sp[0]
		fragment = sp[1]
	}
	return
}

// SubstitutionFunc type used for signature mapping
//
//	key - the section of the url to be replaced: {month} | {month:-2}
//	fragment - would be the -2 in {month:-2}
//	url - the original url
//	d - current time
type SubstitutionFunc func(key string, fragment string, url string, d *time.Time) string

var functions map[string]SubstitutionFunc = map[string]SubstitutionFunc{
	"month":        month,
	"billingMonth": billingMonth,
	"day":          day,
	"billingDay":   billingDay,
}

// month is a SubstitutionFunc type used for signature mapping
//
//	key: the section of the url to be replaced: {month} | {month:-2}
//	fragment: would be the -2 in {month:-2}
//	url: the original url
func month(key string, fragment string, url string, d *time.Time) (m string) {
	m = url
	date := dates.ResetMonth(*d)

	if fragment == "" {
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYM))
		return
	}

	if i, err := strconv.Atoi(fragment); err == nil {
		date = date.AddDate(0, i, 0)
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYM))
	}
	slog.Debug("month",
		slog.String("key", key),
		slog.String("fragment", fragment),
		slog.String("url", url),
		slog.String("m", m),
		slog.String("date", date.String()))
	return
}

func day(key string, fragment string, url string, d *time.Time) (m string) {
	m = url
	date := dates.ResetDay(*d)

	if fragment == "" {
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYMD))
		return
	}

	if i, err := strconv.Atoi(fragment); err == nil {
		date = date.AddDate(0, 0, i)
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYMD))
	}
	slog.Debug("day",
		slog.String("key", key),
		slog.String("fragment", fragment),
		slog.String("url", url),
		slog.String("m", m),
		slog.String("date", date.String()))
	return
}

// billingMonth operates as month, but finds the billing date
func billingMonth(key string, fragment string, url string, d *time.Time) (m string) {
	m = url
	date := dates.BillingEndDate(*d, consts.BILLING_DATE)

	if fragment == "" {
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYM))
		return
	}

	if i, err := strconv.Atoi(fragment); err == nil {
		date = date.AddDate(0, i, 0)
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYM))
	}
	slog.Debug("billingMonth",
		slog.String("key", key),
		slog.String("fragment", fragment),
		slog.String("url", url),
		slog.String("m", m),
		slog.String("date", date.String()))
	return
}

func billingDay(key string, fragment string, url string, d *time.Time) (m string) {
	m = url
	date := dates.BillingEndDate(*d, consts.BILLING_DATE)

	if fragment == "" {
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYMD))
		return
	}

	if i, err := strconv.Atoi(fragment); err == nil {
		date = date.AddDate(0, 0, i)
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYMD))
	}
	slog.Debug("billingMonth",
		slog.String("key", key),
		slog.String("fragment", fragment),
		slog.String("url", url),
		slog.String("m", m),
		slog.String("date", date.String()))
	return
}

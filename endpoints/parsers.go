package endpoints

import (
	"strconv"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/consts"
	"github.com/ministryofjustice/opg-reports/pkg/bi"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
)

// dateArgs uses default values and the passed args slice to work out
// the date modifier, interval period and base date to use for date
// functions
//
// Args:
//
//	0: ago modifier (int)
//	1: base date
func dateArgs(ago int, ts time.Time, args []string) (modifier int, date time.Time) {
	modifier = ago
	date = ts

	if len(args) <= 0 {
		return
	}
	// if args are set, try and work them out
	// 1st arg is present - this should be int representing modification
	if len(args) > 0 {
		if v, err := strconv.Atoi(args[0]); err == nil {
			modifier = v
		}
	}

	// 2nd arg is a base date
	if len(args) > 1 {
		if v, err := time.Parse(consts.DateFormatYearMonthDay, args[1]); err == nil {
			date = v
		}
	}
	return
}

// year provides a YYYY-MM-DD string for the first
// day of the year, with modifications.
// Defaults to -0 years
//
//	{year:-2} => 1st day of the year, 2 months ago
//	{year:-1,2024-02-03} => 2023-01-01
func year(uri string, pg *parserGroup) (u string) {
	u = uri

	ago, date := dateArgs(0, time.Now().UTC(), pg.Arguments)
	date = convert.DateResetYear(date)
	date = convert.DateAddYear(date, ago)

	u = strings.ReplaceAll(uri, pg.Original, date.Format(consts.DateFormatYearMonthDay))

	return
}

// month provides a YYYY-MM-DD string for the first
// day of the month, with modifications.
// Defaults to -9 months
//
//	{month:-2} => 1st day of the month, 2 months ago
//	{month:-1,2024-02-03} => 2024-01-01
func month(uri string, pg *parserGroup) (u string) {
	u = uri

	ago, date := dateArgs(-9, time.Now().UTC(), pg.Arguments)
	date = convert.DateResetMonth(date)
	date = convert.DateAddMonth(date, ago)

	u = strings.ReplaceAll(uri, pg.Original, date.Format(consts.DateFormatYearMonthDay))

	return
}

// day provides a YYYY-MM-DD string for the first
// hour of the day, with modifications.
// Defaults to -1 day
//
//	{day:-2} => 2 days ago at 00:00:00
//	{day:1,2024-02-03} => 2024-02-04
//	{day:0,2024-02-03} => 2024-02-03
func day(uri string, pg *parserGroup) (u string) {
	u = uri

	ago, date := dateArgs(-1, time.Now().UTC(), pg.Arguments)
	date = convert.DateResetDay(date)
	date = convert.DateAddDay(date, ago)

	u = strings.ReplaceAll(uri, pg.Original, date.Format(consts.DateFormatYearMonthDay))

	return
}

// billingMonth provides a YYYY-MM-DD string for the first
// day of the month, base on billing cycle
// Defaults to latest billing month
func billingMonth(uri string, pg *parserGroup) (u string) {
	u = uri
	var date = time.Now().UTC()

	ago, date := dateArgs(0, time.Now().UTC(), pg.Arguments)

	// process the date arguments with defaults
	if date.Day() < consts.CostsBillingDay {
		date = convert.DateAddMonth(date, -2)
	} else {
		date = convert.DateAddMonth(date, -1)
	}

	date = convert.DateResetMonth(date)
	date = convert.DateAddMonth(date, ago)

	u = strings.ReplaceAll(uri, pg.Original, date.Format(consts.DateFormatYearMonthDay))

	return
}

// version updates the url to the current api version
// Any arguments are ignored.
//
//	{version} => v1
//	{version:-1} => v1
func version(uri string, pg *parserGroup) (u string) {
	u = uri
	u = strings.ReplaceAll(u, pg.Original, bi.ApiVersion)
	return
}

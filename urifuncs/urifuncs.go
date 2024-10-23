package urifuncs

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/buildinfo"
	"github.com/ministryofjustice/opg-reports/consts"
	"github.com/ministryofjustice/opg-reports/convert"
)

func dateArgs(mod int, period string, ts time.Time, args ...interface{}) (modifier int, interval string, date time.Time) {
	modifier = mod
	interval = period
	date = ts

	// if args are set, try and work them out
	// 1st arg is present - this should be int representing period modification
	if len(args) > 0 {
		if val, ok := args[0].(int); ok {
			modifier = val
		}
	}
	// 2nd arg is set
	if len(args) > 1 {
		options := []string{"year", "month", "day"}
		if val, ok := args[1].(string); ok && slices.Contains(options, val) {
			interval = val
		}
	}

	// 3nd arg is set
	if len(args) > 2 {
		if val, ok := args[2].(time.Time); ok {
			date = val
		}
	}
	return
}

func dateMod(date time.Time, modifier int, interval string) time.Time {
	// if being asked for a year based start date
	if interval == consts.DateYear {
		date = convert.DateAddYear(convert.DateResetYear(date), modifier)
	} else if interval == consts.DateMonth {
		date = convert.DateAddMonth(convert.DateResetMonth(date), modifier)
	} else if interval == consts.DateDay {
		date = convert.DateAddDay(convert.DateResetDay(date), modifier)
	}
	return date
}

// StartDate replaces {start_date} with a YYYY-MM-DD string
// based on the arguments set
// `args` order:
//  1. period modifier which should be an int (-9, -1 etc)
//  2. interval type - should be year,month,day
//  3. Replacement base date to use
func StartDate(uri string, args ...interface{}) (u string) {
	u = uri
	var (
		match  string    = "{start_date}"
		now    time.Time = time.Now().UTC()
		format string    = consts.DateFormatYearMonthDay
	)
	var (
		modifier int       = -9
		interval string    = "month"
		date     time.Time = now
	)
	// process the date arguments with defaults
	modifier, interval, date = dateArgs(modifier, interval, date, args...)

	date = dateMod(date, modifier, interval)
	dateStr := date.UTC().Format(format)

	slog.Debug("[urifuncs.StartDate]",
		slog.String("uri", uri),
		slog.String("dateStr", dateStr),
		slog.String("interval", interval),
		slog.Int("modifier", modifier),
		slog.String("args", fmt.Sprintf("%+v", args)))

	u = strings.ReplaceAll(u, match, dateStr)

	return
}

// EndDate replaces {end_date} with a YYYY-MM-DD string
// based on the arguments set
// `args` order:
//  1. period modifier which should be an int (-9, -1 etc)
//  2. interval type - should be year,month,day
//  3. Replacement base date to use
func EndDate(uri string, args ...interface{}) (u string) {
	u = uri
	var (
		match  string    = "{end_date}"
		now    time.Time = time.Now().UTC()
		format string    = consts.DateFormatYearMonthDay
	)
	var (
		modifier int       = 0
		interval string    = "month"
		date     time.Time = now
	)
	// process the date arguments with defaults
	modifier, interval, date = dateArgs(modifier, interval, date, args...)

	date = dateMod(date, modifier, interval)
	dateStr := date.UTC().Format(format)

	slog.Debug("[urifuncs.EndDate]",
		slog.String("uri", uri),
		slog.String("dateStr", dateStr),
		slog.String("interval", interval),
		slog.Int("modifier", modifier),
		slog.String("args", fmt.Sprintf("%+v", args)))

	u = strings.ReplaceAll(u, match, dateStr)

	return
}

// BillingEndDate replaces {end_date} with a YYYY-MM-DD string
// of the last billing month
// `args` are unused
func BillingEndDate(uri string, args ...interface{}) (u string) {
	u = uri
	var (
		match  string    = "{end_date}"
		now    time.Time = time.Now().UTC()
		format string    = consts.DateFormatYearMonthDay
		date   time.Time = convert.DateResetMonth(now)
	)

	// process the date arguments with defaults
	if now.Day() < consts.CostsBillingDay {
		date = convert.DateAddMonth(date, -2)
	} else {
		date = convert.DateAddMonth(date, -1)
	}

	dateStr := date.UTC().Format(format)
	slog.Debug("[urifuncs.BillingEndDate]",
		slog.String("uri", uri),
		slog.String("dateStr", dateStr))

	u = strings.ReplaceAll(u, match, dateStr)

	return
}

// BillingStartDate replaces {start_date} with a YYYY-MM-DD string
// of the last billing month
// `args` order:
//  1. period modifier which should be an int (-9, -1 etc)
func BillingStartDate(uri string, args ...interface{}) (u string) {
	u = uri
	var (
		match  string    = "{start_date}"
		now    time.Time = time.Now().UTC()
		format string    = consts.DateFormatYearMonthDay
		date   time.Time = convert.DateResetMonth(now)
	)

	// if theres a modifier, move further back in time
	modifier, _, date := dateArgs(0, "", date, args...)

	// process the date arguments with defaults
	if now.Day() < consts.CostsBillingDay {
		date = convert.DateAddMonth(date, -2)
	} else {
		date = convert.DateAddMonth(date, -1)
	}

	if modifier != 0 {
		date = convert.DateAddMonth(date, modifier)
	}

	dateStr := date.UTC().Format(format)
	slog.Debug("[urifuncs.BillingStartDate]",
		slog.String("uri", uri),
		slog.Int("modifier", modifier),
		slog.String("dateStr", dateStr))

	u = strings.ReplaceAll(u, match, dateStr)

	return
}

// Interval replaces {interval} with suitable value
// `args`:
//  1. alternative interval ("year", "month", "day")
func Interval(uri string, args ...interface{}) (u string) {
	u = uri
	var (
		match    string   = "{interval}"
		interval string   = "month"
		allowed  []string = []string{"year", "month", "day"}
	)

	if len(args) > 0 {
		if val, ok := args[0].(string); ok && slices.Contains(allowed, val) {
			interval = val
		}
	}
	slog.Debug("[urifuncs.Interval]",
		slog.String("uri", uri),
		slog.String("interval", interval))

	u = strings.ReplaceAll(u, match, interval)

	return
}

// Version replaces {version} with the current build version info
// `args` are ignored
func Version(uri string, args ...interface{}) (u string) {
	u = uri
	var match string = "{version}"
	var version string = buildinfo.ApiVersion

	slog.Debug("[urifuncs.Version]",
		slog.String("uri", uri),
		slog.String("version", version))

	u = strings.ReplaceAll(u, match, version)

	return
}

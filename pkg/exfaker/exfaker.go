package exfaker

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/ministryofjustice/opg-reports/pkg/consts"
)

var added bool = false

var (
	FloatMin float64 = -1.5 // Min float value used by float & float_string for random generation.
	FloatMax float64 = 12.5 // Max float value used by float & float_string for random generation.
)

var (
	TimeStringFormat string    = consts.DateFormat                           // DateTime format used in date generation (time_string).
	TimeStringMin    time.Time = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC) // Min time used as lower bound in time_string generation.
	TimeStringMax    time.Time = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC) // Max time used as upper bound in time_string generation.
	DateStringFormat string    = consts.DateFormatYearMonthDay               // Capture the YYYY-MM-DD format
)

// randInt generates a number between min & max
func randInt(min int, max int) int {
	return rand.IntN(max-min) + min
}

// float generates a float64 within the min & max
func float() float64 {
	var (
		min float64 = FloatMin
		max float64 = FloatMax
	)
	return min + rand.Float64()*(max-min)
}

// floatString uses float and then converts to a string version
// using fmt.Sprintf
func floatString() string {
	return fmt.Sprintf("%f", float())
}

// randTime generates a time between min & max values
func randTime() time.Time {
	var max = TimeStringMax
	var min = TimeStringMin
	diff := max.Unix() - min.Unix()
	sec := rand.Int64N(diff) + min.Unix()
	return time.Unix(sec, 0)
}

// timeString returns randTime but as a string using
// the set date format
func timeString() string {
	return randTime().Format(TimeStringFormat)
}

func dateString() string {
	return randTime().Format(DateStringFormat)
}

// uri returns a relative url without hostname etc:
//   - `/test/word/value`
//
// It will create between 1 and 5 path segments
// by calling Word() and appending
func uri() (u string) {
	var n int = randInt(1, 5)
	u = ""

	for i := 0; i < n; i++ {
		var w = faker.Word()
		u += "/" + w
	}

	return
}

// AddProviders to faker for custom versions
// Mostly to generate floats / dates but return them as strings
// for the database models
func AddProviders() {
	slog.Debug("[exfaker.AddProviders] adding")
	if added {
		slog.Debug("[exfaker.AddProviders] already added")
		return
	}

	faker.AddProvider("float", func(v reflect.Value) (interface{}, error) {
		return float(), nil
	})

	faker.AddProvider("float_string", func(v reflect.Value) (interface{}, error) {
		return floatString(), nil
	})

	faker.AddProvider("time_string", func(v reflect.Value) (interface{}, error) {
		return timeString(), nil
	})
	faker.AddProvider("date_string", func(v reflect.Value) (interface{}, error) {
		return dateString(), nil
	})
	faker.AddProvider("uri", func(v reflect.Value) (interface{}, error) {
		return uri(), nil
	})

	added = true
}

// Many returns multiple faked versions of T
func Many[T interface{}](n int, opts ...options.OptionFunc) (faked []*T) {
	slog.Debug("[exfaker.Many] faking many", slog.Int("n", n))

	faked = []*T{}
	for i := 0; i < n; i++ {
		var item T
		var record = &item
		if e := faker.FakeData(record, opts...); e == nil {
			faked = append(faked, record)
		} else {
			slog.Error("[exfaker.Many]", slog.String("err", e.Error()))
		}
	}
	faker.ResetUnique()
	return
}

func init() {
	AddProviders()
}

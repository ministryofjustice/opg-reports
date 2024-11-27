// Package fakerextras provides additional faker methods for go-faker
//
// See the `AddProviders` func for details
package fakerextras

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
)

var added bool = false

const (
	FloatMin         float64 = -1.5             // Min float value used by float & float_string for random generation.
	FloatMax         float64 = 12.5             // Max float value used by float & float_string for random generation.
	TimeStringFormat string  = dateformats.Full // DateTime format used in date generation (time_string).
	DateStringFormat string  = dateformats.YMD  // Capture the YYYY-MM-DD format
)

var (
	now           time.Time = time.Now().UTC()
	TimeStringMin time.Time = now.AddDate(-2, 0, 0) // Min time used as lower bound in time_string generation.
	TimeStringMax time.Time = now.AddDate(0, 0, -5) // Max time used as upper bound in time_string generation.
)

// RandomInt generates a number between min & max
func RandomInt(min int, max int) int {
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
// Bounds are:
//
//	min = TimeStringMin + 2 day
//	max = TimeStringMin - 2 day
//
// This way the time will always be within the dates
// regardless if `>` or `>=` being used
func randTime() time.Time {
	var max = TimeStringMax.AddDate(0, 0, -2)
	var min = TimeStringMin.AddDate(0, 0, 2)
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
	var n int = RandomInt(1, 5)
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
	slog.Debug("[fakerextras] adding")

	if added {
		slog.Debug("[fakerextras] providers already added")
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

type choice interface {
	comparable
}

// Choice will pick a value at random from a list
func Choice[T choice](choices []T) T {
	i := RandomInt(0, len(choices))
	return choices[i]
}

// Choices will pick at least `min` number of items from
// the choices slice passed on
func Choices[T choice](choices []T, min int) (selected []T) {
	var randomised = rand.Perm(len(choices))
	var count = RandomInt(min, len(choices))

	selected = []T{}
	for _, i := range randomised[:count] {
		selected = append(selected, choices[i])
	}

	return
}

// Choose will pick exactly `count` number of items from the choices
func Choose[T choice](choices []T, count int) (selected []T) {
	var randomised = rand.Perm(len(choices))

	selected = []T{}
	for _, i := range randomised[:count] {
		selected = append(selected, choices[i])
	}
	return
}

func init() {
	AddProviders()
}

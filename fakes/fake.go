// Package fake provides some helper functions to generate fake data for mocking / testing
package fakes

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	DateFormat = time.RFC3339
	MaxDate    = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	MinDate    = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
)

// String generates a random string of `length“ from a fixed set of characters
func String(length int) string {
	result := make([]byte, length)
	l := len(charset)
	for i := range result {
		idx := rand.IntN(l)
		result[i] = charset[idx]
	}

	return string(result)
}

// Int generates an int64 whose value is beetween min and max
func Int(min int, max int) int {
	return rand.IntN(max-min) + min
}

// IntAsStr generates a random int (via Int) between min and & max and converts
// that to a string
func IntAsStr(min int, max int) string {
	i := Int(min, max)
	return strconv.Itoa(int(i))
}

// Float creates a float between the min & max
func Float(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// FloatAsStr generates a string version of a randomised float
func FloatAsStr(min float64, max float64) string {
	f := Float(min, max)
	return fmt.Sprintf("%f", f)
}

// Date creates a time.Time between min & max dates
func Date(min time.Time, max time.Time) time.Time {

	diff := max.Unix() - min.Unix()
	sec := rand.Int64N(diff) + min.Unix()

	return time.Unix(sec, 0)
}

// DateAsStr creates a random date between min & max values and retusn string
// version of it
func DateAsStr(min time.Time, max time.Time, f string) string {
	d := Date(min, max)
	return d.Format(f)
}

type IChoice interface {
	comparable
}

func Choice[T IChoice](choices []T) T {
	i := Int(0, len(choices))
	return choices[i]
}

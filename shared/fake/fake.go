package fake

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"
)

// String generates a random string of `length“ from a fixed set of characters
func String(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	l := len(charset)
	for i := range result {
		idx := rand.IntN(l)
		result[i] = charset[idx]
	}

	return string(result)
}

// Int generates an int64 whose value is beetween min and max
func Int(min uint64, max uint64) int64 {
	seed := rand.NewPCG(min, max)
	r := rand.New(seed)
	return r.Int64()
}

// IntAsStr generates a random int (via Int) between min and & max and converts
// that to a string
func IntAsStr(min uint64, max uint64) string {
	i := Int(min, max)
	return strconv.Itoa(int(i))
}

// Float creates a float between the min & max
func Float(min uint64, max uint64) float64 {
	return rand.Float64()
}

// FloatAsStr generates a string version of a randomised float
func FloatAsStr(min uint64, max uint64) string {
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

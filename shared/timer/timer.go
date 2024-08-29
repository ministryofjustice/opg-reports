// Package timer used for timing events
//
// Used by httphandler to time how long calls to urls (like the api)
// take to process
package timer

import (
	"fmt"
	"time"
)

// Ts is for timers in test / benchmarks
type Ts struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Stop sets the value of End to the current time
func (t *Ts) Stop() *Ts {
	t.End = time.Now().UTC()
	return t
}

// Duration works out the diff between the stop and start time
// and returns the seconds gap as a float
func (t *Ts) Duration() float64 {
	if t.End.Year() == 0 {
		t.Stop()
	}
	dur := t.End.Sub(t.Start)
	return dur.Seconds()
}

// Seconds calls Duration and converts that to a string
func (t *Ts) Seconds() string {
	return fmt.Sprintf("%f", t.Duration())
}

func New() *Ts {
	return &Ts{Start: time.Now().UTC()}
}

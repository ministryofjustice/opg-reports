package timer

import (
	"fmt"
	"log/slog"
	"time"
)

// Timer is for timers in test / benchmarks
type Timer struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	duration float64
}

// Stop sets the value of End to the current time
func (t *Timer) Stop() *Timer {
	t.End = time.Now().UTC()
	t.duration = t.End.Sub(t.Start).Seconds()

	slog.Debug("timer duration:", slog.Float64("seconds", t.duration))
	return t
}

// Duration works out the diff between the stop and start time
// and returns the seconds gap as a float
func (t *Timer) Duration() (seconds float64) {
	ranStop := false
	if t.End.Year() == 0 {
		t.Stop()
		ranStop = true
	}
	seconds = t.duration

	if !ranStop {
		slog.Debug("timer duration:", slog.Float64("seconds", seconds))
	}
	return
}

// Seconds calls Duration and converts that to a string
func (t *Timer) Seconds() string {
	return fmt.Sprintf("%f", t.Duration())
}

func New() *Timer {
	return &Timer{Start: time.Now().UTC()}
}

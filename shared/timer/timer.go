// Package timer used to timing events
package timer

import (
	"fmt"
	"time"
)

// Ts is for timers in test / benchmarks
type Ts struct {
	S time.Time `json:"start"`
	E time.Time `json:"end"`
}

func (t *Ts) Stop() *Ts {
	t.E = time.Now().UTC()
	return t
}
func (t *Ts) Duration() float64 {
	if t.E.Year() == 0 {
		t.Stop()
	}
	dur := t.E.Sub(t.S)
	return dur.Seconds()
}

func (t *Ts) Seconds() string {
	return fmt.Sprintf("%f", t.Duration())
}

func New() *Ts {
	return &Ts{S: time.Now().UTC()}
}

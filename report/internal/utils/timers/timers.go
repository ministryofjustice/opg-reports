package timers

import (
	"context"
	"slices"
	"sync"
	"time"
)

type Timer struct {
	Label    string        `json:"label"`
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Dur      time.Duration `json:"-"`
	Duration string        `json:"duration"`
}

type TimerList map[string]*Timer

var (
	mu         sync.Mutex
	contextKey string = "request-timers"
)

// ContextWithTimers will add a timers list to the context for use within all calls
func ContextWithTimers(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, newSet())
}

// Start will trigger a timer for each label passed as with the current time (in UTC)
// as long as the `request-timers` pointer exists within the context.
//
// Use `ContextWithTimers` to get a suitable context object to use with this function,
// as otherwise no timer will be started
//
// Must pass alogn at least 1 label to start a timer, extra labels will start
// additional timers for each label
func Start(ctx context.Context, label string, labels ...string) {
	var list TimerList
	var pntr = getList(ctx)

	if len(labels) == 0 {
		labels = []string{label}
	} else {
		labels = append(labels, label)
	}
	// if no pointer, return immediately
	if pntr == nil {
		return
	}
	list = *pntr
	for _, label := range labels {
		var t = &Timer{Label: label, Start: time.Now().UTC()}
		list[label] = t
	}
	// update the pointer
	mu.Lock()
	pntr = &list
	mu.Unlock()
}

// Stop will set the end time and duration values for timers that match the labels
// passed that are within this context.
//
// If no labels are passed then all timers are stopped.
func Stop(ctx context.Context, labels ...string) {
	var list TimerList
	var pntr = getList(ctx)
	// if no pointer, return immediately
	if pntr == nil {
		return
	}
	list = *pntr
	for k, t := range list {
		var end = time.Now().UTC()
		// if there are no labels or this key is in the label then set any unset durations
		if (len(labels) == 0 || slices.Contains(labels, k)) && t.Duration == "" {
			mu.Lock()
			t.End = end
			t.Dur = end.Sub(t.Start)
			t.Duration = t.Dur.String()
			mu.Unlock()
		}

	}

}

// All returns timers attached to this context
func All(ctx context.Context) (timers []*Timer) {
	var pntr = getList(ctx)
	timers = []*Timer{}

	if pntr == nil {
		return
	}
	for _, t := range *pntr {
		timers = append(timers, t)
	}

	return
}

// getList returns the pointer from the context
func getList(ctx context.Context) *TimerList {
	return ctx.Value(contextKey).(*TimerList)
}

// newSet creates a timerlist to attach to the context
func newSet() *TimerList {
	var set TimerList = map[string]*Timer{}
	return &set
}

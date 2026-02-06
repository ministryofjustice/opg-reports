package timers

import (
	"fmt"
	"sort"
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

var (
	mu  sync.Mutex
	idx int               = 0 // internal default label that gets incremented
	all map[string]*Timer = map[string]*Timer{}
)

// Start will create new times for each label passed along generate start time of now (in UTC). If
// no labels are passed then a counter in used (`t${i}`) to generate one.
//
// Uses a mutex lock
func Start(labels ...string) (list []*Timer) {
	list = []*Timer{}

	// if no label is passed, increasment the base by one
	if len(labels) == 0 {
		idx += 1
		labels = append(labels, fmt.Sprintf("t%d", idx))
	}

	for _, lb := range labels {
		var t = &Timer{Label: lb, Start: time.Now().UTC()}
		mu.Lock()
		all[lb] = t
		list = append(list, t)
		mu.Unlock()
	}

	return
}

// Stop sets the end time and duration details of all timers matching the labels passed along. If
// no labels are passed than it uses the current internal counter to find a timer (`t${i}`).
//
// Uses a mutex lock to update.
func Stop(labels ...string) (list []*Timer) {
	list = []*Timer{}

	if len(labels) == 0 {
		labels = append(labels, fmt.Sprintf("t%d", idx))
	}
	for _, lb := range labels {
		var end = time.Now().UTC()
		// if we find the timer and the duration has not been set, then add end and duration
		if t, ok := all[lb]; ok && t.Duration == "" {
			mu.Lock()
			t.End = end
			t.Dur = end.Sub(t.Start)
			t.Duration = t.Dur.String()
			list = append(list, t)
			mu.Unlock()
		}
	}
	return
}

// AllTimers returns a sorted list of all the currently known timers
func AllTimers() (list []*Timer) {
	list = []*Timer{}
	for _, t := range all {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Dur < list[j].Dur
	})
	return
}

// Clear resets all timers
func Clear() {
	mu.Lock()
	all = map[string]*Timer{}
	mu.Unlock()
}

package timers

import (
	"context"
	"testing"
)

func TestTimers(t *testing.T) {
	var (
		// err error
		test1 context.Context
		test2 context.Context
		l1    []*Timer
		l2    []*Timer
		ctx1  context.Context = context.TODO()
		ctx2  context.Context = context.TODO()
	)

	test1 = ContextWithTimers(ctx1)
	test2 = ContextWithTimers(ctx2)

	Start(test1, "test-001")
	Start(test1, "test-002")
	Start(test2, "test-003")

	Stop(test1, "test-001")
	Stop(test1, "test-002")
	Stop(test2)

	l1 = All(test1)
	if len(l1) != 2 {
		t.Errorf("incorrect number of timers")
	}
	for _, ti := range l1 {
		if ti.Duration == "" {
			t.Errorf("timer should have a duration and be stopped.")
		}
	}

	l2 = All(test2)
	if len(l2) != 1 {
		t.Errorf("inccorect number of timers.")
	}
	for _, ti := range l2 {
		if ti.Duration == "" {
			t.Errorf("timer should have a duration and be stopped.")
		}
	}
}

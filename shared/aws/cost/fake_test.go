package cost

import (
	"testing"
	"time"
)

func TestSharedAwsCostFake(t *testing.T) {
	max := time.Now().UTC()
	min := time.Date(max.Year()-2, max.Month(), 1, 0, 0, 0, 0, time.UTC)
	f := time.RFC3339
	c := Fake(nil, min, max, f)

	d, err := time.Parse(f, c.Date)
	if err != nil {
		t.Errorf("error converting")
	}

	if d.Before(min) || d.After(max) {
		t.Errorf("date setting failed")
	}
	if c.UUID == "" || c.UID() == "" {
		t.Errorf("empty uid")
	}

	c = New(nil)
	fc := Fake(c, min, max, f)
	if fc.UUID != c.UUID {
		t.Errorf("uuid was not reused")
	}

	c = &Cost{}
	fc = Fake(c, min, max, f)
	if fc.UUID == "" {
		t.Errorf("uuid should not be empty")
	}

}

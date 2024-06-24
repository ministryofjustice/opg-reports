package report

import (
	"errors"
	"testing"
)

func TestSharedReportArgs(t *testing.T) {

	a := &ArgNamed{}
	a.SetName("foo")
	if a.GetName() != "foo" {
		t.Errorf("get name failed")
	}

	b := &ArgHelp{}
	b.SetHelp("usage setup")
	if b.GetHelp() != "usage setup" {
		t.Errorf("help / usage failed")
	}

	c := &ArgDefaults{}
	c.SetDefault("1")
	if c.GetDefault() != "1" {
		t.Errorf("default failed: %v", c.GetDefault())
	}

	d := &ArgRequired{}
	d.SetRequired(true)
	if d.GetRequired() != true {
		t.Errorf("required failed")
	}

	arg := NewArg("test", true, "my help", "default")
	if arg.GetName() != "test" {
		t.Errorf("name mismatch")
	}
	if arg.GetHelp() != "my help" {
		t.Errorf("help / usage mis match")
	}
	if arg.GetDefault() != "default" {
		t.Errorf("default failure")
	}
	if arg.GetRequired() != true {
		t.Errorf("required failed")
	}
	val := "input"
	arg.FlagP = &val
	if v, err := arg.Value(); err != nil || v != val {
		t.Errorf("incorrect value")
	}
	if arg.GetFlag() != &val {
		t.Errorf("flag error")
	}

	val = ""
	arg.FlagP = &val
	if _, err := arg.Value(); err == nil || !errors.Is(err, ErrMissingValue) {
		t.Errorf("required value should flag as not set")
	}

	mrg := NewMonthArg("monthtest", true, "my help", "default")
	mval := "2024-04"
	mrg.FlagP = &mval
	if v, err := mrg.Value(); err != nil || v != mval {
		t.Errorf("incorrect value: %v", v)
	}
	if mrg.GetFlag() != &mval {
		t.Errorf("flag error")
	}

}

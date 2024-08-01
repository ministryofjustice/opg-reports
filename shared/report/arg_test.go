package report

import (
	"errors"
	"opg-reports/shared/dates"
	"opg-reports/shared/logger"
	"testing"
	"time"
)

func TestSharedReportArgs(t *testing.T) {
	logger.LogSetup()
	a := &Arg{}
	a.SetName("foo")
	if a.GetName() != "foo" {
		t.Errorf("get name failed")
	}

	b := &Arg{}
	b.SetHelp("usage setup")
	if b.GetHelp() != "usage setup" {
		t.Errorf("help / usage failed")
	}

	c := &Arg{}
	c.SetDefault("1")
	if c.GetDefault() != "1" {
		t.Errorf("default failed: %v", c.GetDefault())
	}

	d := &Arg{}
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
	if v, err := arg.Value(); err != nil || v != val || v != arg.Val() {
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

	mrg := NewMonthArg("monthtest1", true, "my help", "default")
	mval := "2024-04"
	mrg.FlagP = &mval
	if v, err := mrg.Value(); err != nil || v != mval || mrg.Val() != mval {
		t.Errorf("incorrect value: %v", v)
	}
	if mrg.GetFlag() != &mval {
		t.Errorf("flag error")
	}

	mrg = NewMonthArg("monthtest2", true, "my help", "default")
	mval = "not-a-month!"
	mrg.FlagP = &mval
	if _, err := mrg.Value(); err == nil {
		t.Errorf("expected an error for invalid date")
	}

	mrg = NewMonthArg("monthtest3", true, "my help", "default")
	mval = emptyMonth
	now := time.Now().UTC().AddDate(0, -1, 0).Format(dates.FormatYM)
	mrg.FlagP = &mval
	if v, err := mrg.Value(); err != nil || v != now {
		t.Errorf("expected last month for a empty date")
	}

	darg := NewArgConditionalDefault("argdeftest", true, "help", "default!", "null")
	dval := "null"
	darg.FlagP = &dval
	if v, err := darg.Value(); err != nil || v != "default!" {
		t.Errorf("value with a conditional check failed")
	}

	dayArg := NewDayArg("dayarg1", true, "my help", "default")
	dval = emptyDay
	now = time.Now().UTC().AddDate(0, 0, -1).Format(dates.FormatYMD)
	dayArg.FlagP = &dval
	if v, err := dayArg.Value(); err != nil || v != now {
		t.Errorf("expected yesterday for a empty date")
	}

	dayArg = NewDayArg("dayarg2", true, "my help", "default")
	dval = "2024-02-29"
	dayArg.FlagP = &dval
	if v, err := dayArg.Value(); err != nil || v != dval {
		t.Errorf("expected a match")
	}
}

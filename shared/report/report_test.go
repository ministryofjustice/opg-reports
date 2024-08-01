package report

import (
	"opg-reports/shared/logger"
	"testing"
)

func TestSharedReportOK(t *testing.T) {
	logger.LogSetup()
	acc := NewArg("account", false, "account name", "test00")
	month := NewMonthArg("month", true, "month", "")
	in := "2024-01"
	month.FlagP = &in
	args := []IReportArgument{acc, month}

	rep := New()
	rep.SetArguments(args)

	if len(rep.GetArguments()) != len(args) {
		t.Errorf("error getting args")
	}

	if _, err := rep.GetArgument("account"); err != nil {
		t.Errorf("get failed")
	}
	if _, err := rep.GetArgument("not-set"); err == nil {
		t.Errorf("expected error")
	}
	hasRun := false
	rf := func(r IReport) {
		hasRun = true
	}
	rep.SetRunner(rf)
	rep.Run()
	if hasRun != true {
		t.Errorf("did not run")
	}

	fname := rep.Filename()
	if fname != "account^test00.month^2024-01.json" {
		t.Errorf("filename error: %s", fname)
	}

}

func TestSharedReportPanics(t *testing.T) {
	logger.LogSetup()
	defer func() { _ = recover() }()
	month := NewMonthArg("fail", true, "month", "")
	args := []IReportArgument{month}
	rep := New()
	rep.SetArguments(args)
	rep.Run()

	t.Errorf("failed to panic")
}

func TestSharedReportNew(t *testing.T) {
	logger.LogSetup()
	report := New()

	if len(report.GetArguments()) != 0 {
		t.Errorf("length of args error")
	}

	report = New(NewArg("sample", true, "", ""))
	if len(report.GetArguments()) != 1 {
		t.Errorf("length of args error")
	}

}

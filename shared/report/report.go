package report

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"sort"
	"strings"
)

var ErrMissingValue error = errors.New("Required field not set")
var ErrArgNotFound error = errors.New("Argument not found")

// IReportRunF is signature for the function to be run by the command
type IReportRunF func(r IReport)

// IArgs deals with arguments for the command
type IArgs interface {
	SetArguments(arguments []IReportArgument)
	GetArguments() []IReportArgument
	GetArgument(name string) (IReportArgument, error)
}

// IReport is the main runner for a command
type IReport interface {
	SetRunner(runF IReportRunF)
	Run()
	Filename() string
}

// ReportArgs handles getting and setting arguments for the cli
type ReportArgs struct {
	Arguments []IReportArgument
}

// SetArguments overwrites all the arguments
func (r *ReportArgs) SetArguments(arguments []IReportArgument) {
	r.Arguments = arguments
}

// GetArguments returns all the arguments
func (r *ReportArgs) GetArguments() []IReportArgument {
	return r.Arguments
}

// GetArgument returns a single argument matching the name passed
func (r *ReportArgs) GetArgument(name string) (arg IReportArgument, err error) {
	found := false
	for _, a := range r.GetArguments() {
		if a.GetName() == name {
			arg = a
			found = true
		}
	}
	if !found {
		err = ErrArgNotFound
	}
	return
}

// Report is a cmd line report
type Report struct {
	*ReportArgs
	Runner IReportRunF
}

// SetRunner sets the func to call to run
func (r *Report) SetRunner(runF IReportRunF) {
	r.Runner = runF
}

// Run does some pre checks on the arguments then calls the run function
func (r *Report) Run() {
	flag.Parse()
	// Handle validation
	for _, arg := range r.GetArguments() {
		if val, err := arg.Value(); err != nil || val == "" {
			slog.Error("argument error", slog.String("arg", arg.GetName()), slog.String("val", val), slog.String("err", err.Error()))
			panic(err.Error())
		}
	}
	// run the func
	runner := r.Runner
	runner(r)

}

// Filename generates a filename based on the required arguments and their values
func (r *Report) Filename() string {
	str := ""
	mapped := map[string]string{}
	for _, k := range r.GetArguments() {
		mapped[k.GetName()] = k.Val()
	}
	keys := []string{}
	for k := range mapped {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := mapped[key]
		if val == "" {
			val = "-"
		}
		str += fmt.Sprintf("%s^%s.", key, val)
	}
	str = strings.TrimSuffix(str, ".")
	return str + ".json"
}

func New(args ...IReportArgument) *Report {
	rep := &Report{ReportArgs: &ReportArgs{}}
	rep.SetArguments(args)
	return rep
}

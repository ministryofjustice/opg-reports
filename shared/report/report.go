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

type IReportRunF func(r IReport)

type IReport interface {
	SetArguments(arguments []IReportArgument)
	GetArguments() []IReportArgument
	GetArgument(name string) (IReportArgument, error)
	SetRunner(runF IReportRunF)
	Run()
	Filename() string
}

type Report struct {
	Arguments []IReportArgument
	Runner    IReportRunF
}

func (r *Report) SetArguments(arguments []IReportArgument) {
	r.Arguments = arguments
}

func (r *Report) GetArguments() []IReportArgument {
	return r.Arguments
}

func (r *Report) GetArgument(name string) (arg IReportArgument, err error) {
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

func (r *Report) SetRunner(runF IReportRunF) {
	r.Runner = runF
}

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
	rep := &Report{}
	rep.SetArguments(args)
	return rep
}

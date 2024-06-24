package report

import (
	"errors"
	"flag"
	"opg-reports/shared/dates"
	"time"
)

var ErrMonthParse error = errors.New("failed to parse month argument")

type IReportArgumentNamed interface {
	SetName(name string)
	GetName() string
}

type IReportArgumentHelp interface {
	SetHelp(help string)
	GetHelp() string
}

type IReportArgumentDefaults interface {
	SetDefault(def string)
	GetDefault() string
}

type IReportArgumentRequired interface {
	SetRequired(req bool)
	GetRequired() bool
}

type IReportArgument interface {
	IReportArgumentNamed
	IReportArgumentHelp
	IReportArgumentDefaults
	IReportArgumentRequired
	SetFlag()
	GetFlag() *string
	Value() (string, error)
	Val() string
}

type ArgNamed struct {
	Name string
}

func (a *ArgNamed) SetName(name string) {
	a.Name = name
}
func (a *ArgNamed) GetName() string {
	return a.Name
}

type ArgHelp struct {
	Help string
}

func (a *ArgHelp) SetHelp(help string) {
	a.Help = help
}
func (a *ArgHelp) GetHelp() string {
	return a.Help
}

type ArgDefaults struct {
	Default string
}

func (a *ArgDefaults) SetDefault(def string) {
	a.Default = def
}
func (a *ArgDefaults) GetDefault() string {
	return a.Default
}

type ArgRequired struct {
	Required bool
}

func (a *ArgRequired) SetRequired(req bool) {
	a.Required = req
}
func (a *ArgRequired) GetRequired() bool {
	return a.Required
}

type Arg struct {
	ArgNamed
	ArgHelp
	ArgDefaults
	ArgRequired
	FlagP *string
}

func (a *Arg) SetFlag() {
	a.FlagP = flag.String(a.GetName(), a.GetDefault(), a.GetHelp())
}
func (a *Arg) GetFlag() *string {
	return a.FlagP
}

func (a *Arg) Value() (val string, err error) {
	value := *a.FlagP
	if value != "" || len(value) > 0 {
		val = value
	} else if a.Required {
		err = ErrMissingValue
	}
	return
}

func (a *Arg) Val() string {
	v, _ := a.Value()
	return v
}

type MonthArg struct {
	ArgNamed
	ArgHelp
	ArgDefaults
	ArgRequired
	FlagP *string
}

func (a *MonthArg) SetFlag() {
	a.FlagP = flag.String(a.GetName(), a.GetDefault(), a.GetHelp())
}
func (a *MonthArg) GetFlag() *string {
	return a.FlagP
}

func (a *MonthArg) Value() (val string, err error) {
	value, e := a.MonthValue()

	if e != nil {
		err = e
	} else if a.GetRequired() && value.Format(dates.FormatY) == dates.ErrYear {
		err = ErrMonthParse
	} else {
		val = value.Format(dates.FormatYM)
	}
	return
}
func (a *MonthArg) Val() string {
	v, _ := a.Value()
	return v
}

func (a *MonthArg) MonthValue() (val time.Time, err error) {
	return dates.StringToDate(*a.FlagP)
}

func NewArg(name string, required bool, usage string, def string) *Arg {
	arg := &Arg{
		ArgNamed:    ArgNamed{Name: name},
		ArgHelp:     ArgHelp{Help: usage},
		ArgDefaults: ArgDefaults{Default: def},
		ArgRequired: ArgRequired{Required: required},
	}
	arg.SetFlag()
	return arg
}

func NewMonthArg(name string, required bool, usage string, def string) *MonthArg {
	arg := &MonthArg{
		ArgNamed:    ArgNamed{Name: name},
		ArgHelp:     ArgHelp{Help: usage},
		ArgDefaults: ArgDefaults{Default: def},
		ArgRequired: ArgRequired{Required: required},
	}
	arg.SetFlag()
	return arg
}

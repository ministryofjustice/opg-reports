package report

import (
	"errors"
	"flag"
	"opg-reports/shared/dates"
	"time"
)

var ErrMonthParse error = errors.New("failed to parse month argument")
var ErrDayParse error = errors.New("failed to parse day argument")

// IReportArgumentNamed handles getting and setting the name of argument,
// this will be the cli flag as well
type IReportArgumentNamed interface {
	SetName(name string)
	GetName() string
}

// IReportArgumentHelp deals with setting and getting arg usage / help
// messages
type IReportArgumentHelp interface {
	SetHelp(help string)
	GetHelp() string
}

// IReportArgumentDefaults handles default value for an argument
type IReportArgumentDefaults interface {
	SetDefault(def string)
	GetDefault() string
}

// IReportArgumentRequired handles setting if this argument is required
// or not. Required arguments get checked for a value being present
// by the IReport and are used within filename generation
type IReportArgumentRequired interface {
	SetRequired(req bool)
	GetRequired() bool
}

// IReport flag handles creating the argument using something
// like the flag package
type IReportFlag interface {
	SetFlag()
	GetFlag() *string
}

// IReportArgument models a command line argument
type IReportArgument interface {
	IReportArgumentNamed
	IReportArgumentHelp
	IReportArgumentDefaults
	IReportArgumentRequired
	IReportFlag
	Value() (string, error)
	Val() string
}

// Arg is a standard argument for a report
type Arg struct {
	Name             string
	Help             string
	Default          string
	DefaultCondition *string
	Required         bool
	FlagP            *string
}

// SetName replaces the name
func (a *Arg) SetName(name string) {
	a.Name = name
}

// GetName returns the name
func (a *Arg) GetName() string {
	return a.Name
}

// SetHelp replaces the help message
func (a *Arg) SetHelp(help string) {
	a.Help = help
}

// GetHelp returns the usage message
func (a *Arg) GetHelp() string {
	return a.Help
}

// SetDefault replaces the default value
func (a *Arg) SetDefault(def string) {
	a.Default = def
}

// GetDefault returns the value
func (a *Arg) GetDefault() string {
	return a.Default
}

// SetRequired replaces the value
func (a *Arg) SetRequired(req bool) {
	a.Required = req
}

// GetRequired returns the value
func (a *Arg) GetRequired() bool {
	return a.Required
}
func (a *Arg) SetFlag() {
	a.FlagP = flag.String(a.GetName(), a.GetDefault(), a.GetHelp())
}
func (a *Arg) GetFlag() *string {
	return a.FlagP
}

// Value uses the flag value and returns string version
// if theres an error, will return that
func (a *Arg) Value() (val string, err error) {
	value := *a.FlagP
	defCond := a.DefaultCondition
	if value != "" && defCond != nil && value == *defCond {
		val = a.Default
	} else if value != "" || len(value) > 0 {
		val = value
	} else if a.Required {
		err = ErrMissingValue
	}
	return
}

// Val call Value(), but disregards the error and in that case returns
// empty string
func (a *Arg) Val() string {
	v, _ := a.Value()
	return v
}

const emptyMonth string = "-"
const emptyDay string = "-"

// MonthArg is a custom arg that represents a YYYY-MM inputed value
// and replaces the Value() to handle parsing of string to time.Time
type MonthArg struct {
	*Arg
}

// Value fetches the value and tries to parse this into
// and time.Time (via MonthValue func)
// If the parsing fails or if the value matches a default date
// (0000-01) then return an error message, otherwise
// return YYYY-MM version of the inputed date
func (a *MonthArg) Value() (val string, err error) {
	value, e := a.MonthValue()
	if e != nil {
		err = e
		return
	} else if a.GetRequired() && value.Format(dates.FormatY) == dates.ErrYear {
		err = ErrMonthParse
	} else {
		val = value.Format(dates.FormatYM)
	}
	return
}

// Val calls Value(), but disregards the error message
func (a *MonthArg) Val() string {
	v, _ := a.Value()
	return v
}

// MonthValue is used to convert the string version of the argument
// into a time.Time
func (a *MonthArg) MonthValue() (val time.Time, err error) {
	rawValue := *a.FlagP
	if rawValue == emptyMonth {
		val = time.Now().UTC().AddDate(0, -1, 0)
	} else {
		val, err = dates.StringToDate(*a.FlagP)
	}
	return
}

// DayArg is a custom arg that represents a YYYY-MM-DD inputed value
// and replaces the Value() to handle parsing of string to time.Time
type DayArg struct {
	*Arg
}

// Value fetches the value and tries to parse this into
// and time.Time (via DayValue func)
// If the parsing fails or if the value matches a default date
// (0000-01) then return an error message, otherwise
// return YYYY-MM-DD version of the inputed date
func (a *DayArg) Value() (val string, err error) {
	value, e := a.DayValue()
	if e != nil {
		err = e
		return
	} else if a.GetRequired() && value.Format(dates.FormatY) == dates.ErrYear {
		err = ErrDayParse
	} else {
		val = value.Format(dates.FormatYMD)
	}
	return
}

// Val calls Value(), but disregards the error message
func (a *DayArg) Val() string {
	v, _ := a.Value()
	return v
}

// DayValue is used to convert the string version of the argument
// into a time.Time
func (a *DayArg) DayValue() (val time.Time, err error) {
	rawValue := *a.FlagP
	if rawValue == emptyDay {
		n := time.Now().UTC().AddDate(0, 0, -1)
		val = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		val, err = dates.StringToDate(*a.FlagP)
	}
	return
}

// NewArg generates a new argument
func NewArg(name string, required bool, usage string, def string) *Arg {

	arg := &Arg{
		Name:             name,
		Help:             usage,
		Default:          def,
		DefaultCondition: nil,
		Required:         required,
	}
	arg.SetFlag()
	return arg
}

func NewArgConditionalDefault(name string, required bool, usage string, def string, condVal string) *Arg {

	arg := &Arg{
		Name:             name,
		Help:             usage,
		Default:          def,
		DefaultCondition: &condVal,
		Required:         required,
	}
	arg.SetFlag()
	return arg
}

// NewMonthArg generates a new month argument
func NewMonthArg(name string, required bool, usage string, def string) *MonthArg {
	a := NewArg(name, required, usage, def)
	arg := &MonthArg{
		Arg: a,
	}

	return arg
}

// NewDayArg generates a new month argument
func NewDayArg(name string, required bool, usage string, def string) *DayArg {
	a := NewArg(name, required, usage, def)
	arg := &DayArg{
		Arg: a,
	}
	return arg
}

package report

import (
	"errors"
	"flag"
	"opg-reports/shared/dates"
	"time"
)

var ErrMonthParse error = errors.New("failed to parse month argument")

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

// ArgNamed handles tracking the name of the argument
type ArgNamed struct {
	Name string
}

// SetName replaces the name
func (a *ArgNamed) SetName(name string) {
	a.Name = name
}

// GetName returns the name
func (a *ArgNamed) GetName() string {
	return a.Name
}

// ArgHelp handles usage message
type ArgHelp struct {
	Help string
}

// SetHelp replaces the help message
func (a *ArgHelp) SetHelp(help string) {
	a.Help = help
}

// GetHelp returns the usage message
func (a *ArgHelp) GetHelp() string {
	return a.Help
}

// ArgDefaults handles default values
type ArgDefaults struct {
	Default string
}

// SetDefault replaces the default value
func (a *ArgDefaults) SetDefault(def string) {
	a.Default = def
}

// GetDefault returns the value
func (a *ArgDefaults) GetDefault() string {
	return a.Default
}

// ArgRequired determines if its required or not
type ArgRequired struct {
	Required bool
}

// SetRequired replaces the value
func (a *ArgRequired) SetRequired(req bool) {
	a.Required = req
}

// GetRequired returns the value
func (a *ArgRequired) GetRequired() bool {
	return a.Required
}

// ArgFlag handles the actual flag package usage
type ArgFlag struct {
	*ArgNamed
	*ArgHelp
	*ArgDefaults
	FlagP *string
}

func (a *ArgFlag) SetFlag() {
	a.FlagP = flag.String(a.GetName(), a.GetDefault(), a.GetHelp())
}
func (a *ArgFlag) GetFlag() *string {
	return a.FlagP
}

// Arg is a standard argument for a report
type Arg struct {
	*ArgRequired
	*ArgFlag
}

// Value uses the flag value and returns string version
// if theres an error, will return that
func (a *Arg) Value() (val string, err error) {
	value := *a.FlagP
	if value != "" || len(value) > 0 {
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

// MonthArg is a custom arg that represents a YYYY-MM inputed value
// and replaces the Value() to handle parsing of string to time.Time
type MonthArg struct {
	*ArgRequired
	*ArgFlag
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
	return dates.StringToDate(*a.FlagP)
}

// NewArg generates a new argument
func NewArg(name string, required bool, usage string, def string) *Arg {
	argFlag := &ArgFlag{
		ArgNamed:    &ArgNamed{Name: name},
		ArgHelp:     &ArgHelp{Help: usage},
		ArgDefaults: &ArgDefaults{Default: def},
	}
	arg := &Arg{
		ArgFlag:     argFlag,
		ArgRequired: &ArgRequired{Required: required},
	}
	arg.SetFlag()
	return arg
}

// NewMonthArg generates a new month argument
func NewMonthArg(name string, required bool, usage string, def string) *MonthArg {
	argFlag := &ArgFlag{
		ArgNamed:    &ArgNamed{Name: name},
		ArgHelp:     &ArgHelp{Help: usage},
		ArgDefaults: &ArgDefaults{Default: def},
	}
	arg := &MonthArg{
		ArgFlag:     argFlag,
		ArgRequired: &ArgRequired{Required: required},
	}
	arg.SetFlag()
	return arg
}

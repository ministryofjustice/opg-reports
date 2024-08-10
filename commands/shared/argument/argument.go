package argument

import (
	"flag"
	"strconv"
	"time"
)

const blank string = "-"

// -- Str
type Str struct {
	Value   *string
	Default string
}

func (s Str) String() (str string) {
	if s.Value != nil {
		str = *s.Value
	}
	return
}
func (s Str) Set(str string) (err error) {
	*s.Value = str
	return
}

type Int struct {
	Value   *int
	Default int
}

func (i Int) String() (str string) {
	if i.Value != nil {
		str = strconv.Itoa(*i.Value)
	}
	return
}
func (i Int) Set(str string) (err error) {
	if x, err := strconv.Atoi(str); err == nil {
		*i.Value = x
	}
	return
}

// -- Date
type Date struct {
	Value   *time.Time
	Default time.Time
	format  string
}

func (d Date) String() (str string) {
	if d.Value != nil {
		str = d.Value.Format(d.format)
	}
	return
}

func (d Date) Set(str string) (err error) {
	if t, err := time.Parse(d.format, str); err == nil {
		*d.Value = t
	}
	return
}

func New(fs *flag.FlagSet, name string, def string, usage string) *Str {
	var s = def
	var arg = &Str{Value: &s, Default: def}
	fs.Var(arg, name, usage)
	return arg
}

func NewInt(fs *flag.FlagSet, name string, def int, usage string) *Int {
	var s = def
	var arg = &Int{Value: &s, Default: def}
	fs.Var(arg, name, usage)
	return arg
}

func NewDate(fs *flag.FlagSet, name string, def time.Time, format string, usage string) *Date {
	var t = def
	var arg = &Date{Value: &t, format: format, Default: def}
	fs.Var(arg, name, usage)
	return arg
}

package response

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"time"
)

// ITimings handles simple start, end and duration elements of the interface.
type ITimings interface {
	Start()
	GetStart() time.Time
	End()
	GetEnd() time.Time
	GetDuration() time.Duration
}

type IDataTimings interface {
	AddTimestamp(ts time.Time)
	GetMinMax() (min time.Time, max time.Time)
}

type ITimingData interface {
	ITimings
	IDataTimings
}

// IStatus handles tracking the http status of the api response.
// Its value should be used with IApi.Write call at the end
type IStatus interface {
	SetStatus(status int)
	GetStatus() int
}

// IErrors allows tracking of server side errors such as validation
// and will be included in the IApi.Write
type IErrors interface {
	SetErrors(errors []error)
	AddError(err error)
	GetErrors() []error
}

// IBase is a merge interface that wuld be typical of an http response.
// This version excludes any result data / handling for simplicty on errors or
// empty results
type IBase interface {
	ITimingData
	IStatus
	IErrors
	AddErrorWithStatus(err error, status int)
}

// IResult providers a response interface whose result type can vary between
// slice, a map or a map of slices.
// This allows api respsones to adapt to the most useful data type for the endpoint
type IResult[C ICell, R IRow[C], D ITableData[C, R]] interface {
	IBase
	SetResult(result D)
	GetResult() D
	GetDataTimings() (min *time.Time, max *time.Time)
}

type dataTimings struct {
	Min *time.Time `json:"min"`
	Max *time.Time `json:"max"`
	All []int64    `json:"-"`
}
type requestTimes struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

// Timings impliments [ITimings] & [IDataTimings] which is [ITimingData]
type Timings struct {
	RequestTimes *requestTimes `json:"timings"`
	Datatimes    *dataTimings  `json:"data_timestamps,omitempty"`
}

// Start tracks the start time of this request
func (i *Timings) Start() {
	i.RequestTimes.Start = time.Now().UTC()
}

// End tracks the end time and the duration of the request
func (i *Timings) End() {
	i.RequestTimes.End = time.Now().UTC()
	i.RequestTimes.Duration = i.RequestTimes.End.Sub(i.RequestTimes.Start)
}
func (i *Timings) GetStart() time.Time {
	return i.RequestTimes.Start
}

func (i *Timings) GetEnd() time.Time {
	return i.RequestTimes.End
}
func (i *Timings) GetDuration() time.Duration {
	return i.RequestTimes.Duration
}

func (i *Timings) AddTimestamp(ts time.Time) {
	u := ts.UnixMicro()
	i.Datatimes.All = append(i.Datatimes.All, u)
	// check the min max values
	if i.Datatimes.Max == nil || i.Datatimes.Min == nil {
		i.GetMinMax()
	} else if m := i.Datatimes.Min.Unix(); u < m {
		i.GetMinMax()
	} else if m := i.Datatimes.Min.Unix(); u > m {
		i.GetMinMax()
	}
}
func (i *Timings) GetMinMax() (minT time.Time, maxT time.Time) {
	uMin := slices.Min(i.Datatimes.All)
	uMax := slices.Max(i.Datatimes.All)

	minT = time.UnixMicro(uMin).UTC()
	maxT = time.UnixMicro(uMax).UTC()

	i.Datatimes.Min = &minT
	i.Datatimes.Max = &maxT
	return
}

// Status impliments [IStatus]
// Provides http status tracking
type Status struct {
	Code int `json:"status"`
}

// SetStatus updates the status field
func (i *Status) SetStatus(status int) {
	i.Code = status
}

// GetStatus returns the status field
func (i *Status) GetStatus() int {
	return i.Code
}

// Errors handles error tracking and impliments [IErrors]
type Errors struct {
	Errs []error `json:"errors"`
}

// SetErrors replaces errors with those passed
func (r *Errors) SetErrors(errors []error) {
	r.Errs = errors
}

// AddError add a new error to the list
func (r *Errors) AddError(err error) {
	r.Errs = append(r.Errs, err)
}

// GetErrors returns all errors
func (r *Errors) GetErrors() []error {
	return r.Errs
}

// Base impliments IBase
// Would be used for a simple endpoint that doesn't return data,
// such as an api root
type Base struct {
	*Timings
	*Status
	*Errors
}

// AddErrorWithStatus adds an error and updates the status at the same time.
// Helpful when validating fields to do both at once.
func (i *Base) AddErrorWithStatus(err error, status int) {
	i.AddError(err)
	i.SetStatus(status)
}

// Result impliments [IResult].
// It allows a response to return with variable (C) data type. This is currently
// constrained to map[string]R, map[string][]R and []R.
// This means various enpoints can return differing ways collecting the data.IEntry
// so some can group by a field or just list everything that matches
//
// This struct and interface allows you to easily decode a response as long as you know
// its return type
type Result[C ICell, R IRow[C], D ITableData[C, R]] struct {
	Base
	Res D `json:"result"`
}

// SetResult updates the internal result data
func (i *Result[C, R, D]) SetResult(result D) {
	i.Res = result
}

// GetResult returns the result
func (i *Result[C, R, D]) GetResult() D {
	return i.Res
}

// GetDataTimings returns the rough range of the age of the data included in the response
func (i *Result[C, R, D]) GetDataTimings() (min *time.Time, max *time.Time) {
	if i.Timings.Datatimes != nil {
		min = i.Timings.Datatimes.Min
		max = i.Timings.Datatimes.Max
	}
	return
}

// NewSimpleResult returns a fresh Base with
// status set as OK and errors empty
func NewSimpleResult() *Base {
	return &Base{
		Timings: &Timings{Datatimes: &dataTimings{}, RequestTimes: &requestTimes{}},
		Status:  &Status{Code: http.StatusOK},
		Errors:  &Errors{Errs: []error{}},
	}
}

func NewResponse() *Result[*Cell, *Row[*Cell], *TableData[*Cell, *Row[*Cell]]] {
	return &Result[*Cell, *Row[*Cell], *TableData[*Cell, *Row[*Cell]]]{
		Base: *NewSimpleResult(),
	}
}

func ParseFromJson[C ICell, R IRow[C], D ITableData[C, R]](content []byte, response *Result[C, R, D]) (err error) {
	err = json.Unmarshal(content, response)
	return
}
func ParseFromHttp[C ICell, R IRow[C], D ITableData[C, R]](r *http.Response, response *Result[C, R, D]) (err error) {
	_, by := Stringify(r)
	return ParseFromJson(by, response)
}

// Stringify takes a http.Response and returns string & []byte
// values of the response body
func Stringify(r *http.Response) (string, []byte) {
	b, _ := io.ReadAll(r.Body)
	return string(b), b
}
package response

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"time"
)

// IRequestStart track the start time for the request
type IRequestStart interface {
	SetStart()
	GetStart() time.Time
}

// IRequestEnd track the end time of the response
type IRequestEnd interface {
	SetEnd()
	GetEnd() time.Time
}

// IRequestDuration handles the duration of the request
type IRequestDuration interface {
	IRequestStart
	IRequestEnd
	SetDuration()
	GetDuration() time.Duration
}

// IDataRecency tracks the data creation times to inform the
// age of the data
type IDataRecency interface {
	SetDataAge(ts ...time.Time)
	GetDataAgeMin() time.Time
	GetDataAgeMax() time.Time
}

// IResponseStatus handles tracking the http status of the api response.
// Its value should be used with IApi.Write call at the end
type IResponseStatus interface {
	SetStatus(status int)
	GetStatus() int
}

// IErrors allows tracking of server side errors such as validation
// and will be included in the IApi.Write
type IErrors interface {
	SetError(errors ...error)
	GetError() []error
}

// IErrorWithStatus allows adding an error message and changing the response
// status at the same time
type IErrorWithStatus interface {
	IResponseStatus
	IErrors
	SetErrorAndStatus(err error, status int)
}

type IResponseData[C ICell, R IRow[C]] interface {
	SetData(t ITable[C, R])
	GetData() ITable[C, R]
}

type IResponse[C ICell, R IRow[C]] interface {
	IRequestStart
	IRequestEnd
	IRequestDuration
	IDataRecency
	IResponseStatus
	IErrors
	IErrorWithStatus
	IResponseData[C, R]
}

// --- CONCREATE VERSIONS

type requestTimes struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}
type dataAge struct {
	Min *time.Time `json:"min"`
	Max *time.Time `json:"max"`
	All []int64    `json:"-"`
}
type Response[C ICell, R IRow[C]] struct {
	RequestTimes *requestTimes `json:"request_timings,omitempty"`
	DataAge      *dataAge      `json:"data_age"`
	StatusCode   int           `json:"status"`
	Errors       []error       `json:"errors"`
	Data         ITable[C, R]  `json:"result"`
}

// --- IResponseData

func (r *Response[C, R]) SetData(t ITable[C, R]) {
	r.Data = t
}
func (r *Response[C, R]) GetData() ITable[C, R] {
	return r.Data
}

// --- IErrorWithStatus

func (r *Response[C, R]) SetErrorAndStatus(err error, status int) {
	r.SetError(err)
	r.SetStatus(status)
}

// --- IErrors

func (r *Response[C, R]) SetError(errors ...error) {
	if errors == nil {
		r.Errors = []error{}
	} else {
		r.Errors = append(r.Errors, errors...)
	}
}

func (r *Response[C, R]) GetError() []error {
	return r.Errors
}

// --- IResponseStatus

func (r *Response[C, R]) SetStatus(status int) {
	r.StatusCode = status
}
func (r *Response[C, R]) GetStatus() int {
	return r.StatusCode
}

// --- IDataRecency

func (r *Response[C, R]) SetDataAge(times ...time.Time) {
	if times == nil {
		r.DataAge.All = []int64{}
	} else {
		for _, t := range times {
			r.DataAge.All = append(r.DataAge.All, t.UnixMicro())
		}
	}
}

func (r *Response[C, R]) GetDataAgeMin() (t time.Time) {
	if r.DataAge.Min != nil {
		return *r.DataAge.Min
	} else if len(r.DataAge.All) > 0 {
		min := slices.Min(r.DataAge.All)
		t = time.UnixMicro(min).UTC()
		r.DataAge.Min = &t
	}
	return t
}

func (r *Response[C, R]) GetDataAgeMax() (t time.Time) {
	if r.DataAge.Max != nil {
		return *r.DataAge.Max
	} else if len(r.DataAge.All) > 0 {
		max := slices.Max(r.DataAge.All)
		t = time.UnixMicro(max).UTC()
		r.DataAge.Max = &t
	}
	return t
}

// --- IRequestStart

func (r *Response[C, R]) SetStart() {
	r.RequestTimes.Start = time.Now().UTC()
}
func (r *Response[C, R]) GetStart() time.Time {
	return r.RequestTimes.Start
}

// --- IRequestEnd

func (r *Response[C, R]) SetEnd() {
	r.RequestTimes.End = time.Now().UTC()
}
func (r *Response[C, R]) GetEnd() time.Time {
	return r.RequestTimes.End
}

// -- IRequestDuration

func (r *Response[C, R]) SetDuration() {
	r.RequestTimes.Duration = r.GetEnd().Sub(r.GetStart())
}
func (r *Response[C, R]) GetDuration() time.Duration {
	return r.RequestTimes.Duration
}

// --- NEW HELPERS

func NewResponse[C ICell, R IRow[C]]() IResponse[C, R] {
	return &Response[C, R]{
		RequestTimes: &requestTimes{},
		DataAge:      &dataAge{},
		StatusCode:   http.StatusOK,
		Errors:       []error{},
		Data:         NewTable[C, R](),
	}
}

func ToJson[C ICell, R IRow[C]](r IResponse[C, R]) (content []byte, err error) {
	return json.MarshalIndent(r, "", "  ")
}

func FromJson[C ICell, R IRow[C]](content []byte, r IResponse[C, R]) (err error) {
	err = json.Unmarshal(content, r)
	return
}

func FromHttp[C ICell, R IRow[C]](content *http.Response, r IResponse[C, R]) (err error) {
	_, bytes := Stringify(content)
	return FromJson[C, R](bytes, r)
}

func Stringify(r *http.Response) (s string, b []byte) {
	b, _ = io.ReadAll(r.Body)
	s = string(b)
	return
}

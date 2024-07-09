// Package response provides interfaces and structs for an IApi response object.
//
// The IApi utilises an IRepsonse object to provide the content that will be
// sent back in the request
//
// By default, we want that resposne to include start, end and duration
// of the original request - this will help and tracking and performance
// analysis of the api
//
// We also want to ensure we include a http status code as well errors
// about the reason the request failed (validation etc).
//
// For more complex api calls we also provide metadata about the original
// request - this should contain things like query conditions found within
// the request (path values or query strings) so this can be confirmed
// on the recievers end
//
// The main data included in the response is modelled as a table and
// handled by ITable, which is also within this package
//
// The response data (GetData, SetData with ITable) is s single table,
// analogus to a spreadsheet / html table. It contains head, body and foot
// elements with getters & setters for each
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

// IResponseDataRecency tracks the data creation times to inform the
// age of the data
type IResponseDataRecency interface {
	SetDataAge(ts ...time.Time)
	// GetAllDataAge() []time.Time
	GetDataAgeMin() time.Time
	GetDataAgeMax() time.Time
}

// IResponseStatus handles tracking the http status of the api response.
// Its value should be used with IApi.Write call at the end
type IResponseStatus interface {
	SetStatus(status int)
	GetStatus() int
}

// IResponseErrors allows tracking of server side errors such as validation
// and will be included in the IApi.Write
type IResponseErrors interface {
	SetError(errors ...error)
	GetError() []error
}

// IResponseErrorWithStatus allows adding an error message and changing the response
// status at the same time
type IResponseErrorWithStatus interface {
	IResponseStatus
	IResponseErrors
	SetErrorAndStatus(err error, status int)
}

type IResponseMetadata interface {
	SetMetadata(k string, v interface{})
	GetMetadata() map[string]interface{}
}

type IResponseData[C ICell, R IRow[C]] interface {
	SetData(t ITable[C, R])
	GetData() ITable[C, R]
}

type IResponse[C ICell, R IRow[C]] interface {
	IRequestStart
	IRequestEnd
	IRequestDuration
	IResponseDataRecency
	IResponseStatus
	IResponseErrors
	IResponseErrorWithStatus
	IResponseMetadata
	IResponseData[C, R]
}

// --- CONCREATE VERSIONS
// requestTimes handles start, end & duration of the request
type requestTimes struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

// dataAge is used to track the age of the data items in the response
// so we can display a "data accurate as of" message if wanted
type dataAge struct {
	Min *time.Time `json:"min"`
	Max *time.Time `json:"max"`
	All []int64    `json:"-"`
}

// Response is the main response struct that is returned by API
// endpoints.
//
// All data is stored and ahndles as a table, made up of header, body
// and foot liks a HTML table (as that is main how the data is used).
//
// Errors and Http status codes are also tracked
// Impliments [IResponse]
type Response[C ICell, R IRow[C]] struct {
	RequestTimes *requestTimes          `json:"request_timings,omitempty"`
	DataAge      *dataAge               `json:"data_age"`
	StatusCode   int                    `json:"status"`
	Errors       []error                `json:"errors"`
	Metadata     map[string]interface{} `json:"metadata"`
	Data         ITable[C, R]           `json:"result"`
}

// --- IResponseData
// SetData overwrites the main data object of the repsonse with a new
// table of data.
//
// Called towards the end of a http handler to set the response ready
// for transmission
//
// Interface: [IResponseData]
func (r *Response[C, R]) SetData(t ITable[C, R]) {
	r.Data = t
}

// GetData returns the tabular data from this API response
//
// Interface: [IResponseData]
func (r *Response[C, R]) GetData() ITable[C, R] {
	return r.Data
}

// Set response metadata, this can include relevant info
// such as query string values and filter requests
//
// Interface: [IResponseMetadata]
func (r *Response[C, R]) SetMetadata(k string, v interface{}) {
	r.Metadata[k] = v
}

// GetMetadata
//
// Interface: [IResponseMetadata]
func (r *Response[C, R]) GetMetadata() map[string]interface{} {
	return r.Metadata
}

// --- IResponseErrorWithStatus
// SetErrorAndStatus adds an error to the stack of erroes and also
// changes the http status code at the same time. This provides
// slightly cleaner error handling with http func by combining two calls
//
// Interface: [IResponseErrorWithStatus]
func (r *Response[C, R]) SetErrorAndStatus(err error, status int) {
	r.SetError(err)
	r.SetStatus(status)
}

// --- IResponseErrors
// SetError appends a new error the error stack. This is then
// included within the response from the api so issues with
// the request can be debugged.
//
// Note: If `nil` is passed, then the error stack is reset
// Interface: [IResponseErrors]
func (r *Response[C, R]) SetError(errors ...error) {
	if errors == nil {
		r.Errors = []error{}
	} else {
		r.Errors = append(r.Errors, errors...)
	}
}

// GetError returns all stored errors in this response
//
// Interface: [IResponseErrors]
func (r *Response[C, R]) GetError() []error {
	return r.Errors
}

// --- IResponseStatus
// SetStatus updates the http status code that will be used
// within the respone to the value passed. When use NewResponse
// this is set as http.StatusOK by default.
//
// Interface: [IResponseStatus]
func (r *Response[C, R]) SetStatus(status int) {
	r.StatusCode = status
}

// GetStatus returns the current status code value
//
// Interface: [IResponseStatus]
func (r *Response[C, R]) GetStatus() int {
	return r.StatusCode
}

// --- IResponseDataRecency
// SetDataAge is used to track the created times of the data
// associated with this api response. This provides a way
// to show messages about when the data was created / updated
// in front ends
//
// Interface: [IResponseDataRecency]
func (r *Response[C, R]) SetDataAge(times ...time.Time) {
	if times == nil {
		r.DataAge.All = []int64{}
	} else {
		for _, t := range times {
			r.DataAge.All = append(r.DataAge.All, t.UnixMicro())
		}
	}
}

// GetDataAgeMin returns the min date stored within the
// data in this api response (effectively the "oldest")
//
// Interface: [IResponseDataRecency]
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

// GetDataAgeMax returns the max date stored within the
// data in this api response (effectively the "youngest")
//
// Interface: [IResponseDataRecency]
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

// SetStart is callef early in the http request handler to
// accurately track the start time of the request.
// This is then used with end time to work out durations
// that can be checked for performance etc
//
// Interface: [IRequestStart]
func (r *Response[C, R]) SetStart() {
	r.RequestTimes.Start = time.Now().UTC()
}

// GetStart returns the request start time data
// Interface: [IRequestStart]
func (r *Response[C, R]) GetStart() time.Time {
	return r.RequestTimes.Start
}

// --- IRequestEnd

// SetEnd is companiion to SetStart, this tracks the
// end of the http request being processed and can then
// be used to work out duration.
//
// Interface: [IRequestEnd]
func (r *Response[C, R]) SetEnd() {
	r.RequestTimes.End = time.Now().UTC()
}

// GetEnd returns the end time of the request
//
// Interface: [IRequestEnd]
func (r *Response[C, R]) GetEnd() time.Time {
	return r.RequestTimes.End
}

// -- IRequestDuration

// SetDuration uses the end & start times to work out how long
// a http request has taken to be processed. This informaton
// is helpful for assessing performance over the api
//
// Interface: [IRequestDuration]
func (r *Response[C, R]) SetDuration() {
	r.RequestTimes.Duration = r.GetEnd().Sub(r.GetStart())
}

// GetDuration returns the duration data
//
// Interface: [IRequestDuration]
func (r *Response[C, R]) GetDuration() time.Duration {
	return r.RequestTimes.Duration
}

// --- NEW HELPERS

// NewResponse returns a fresh response with empty data setup
// Note: StatusCode is set to http.StatusOk
func NewResponse[C ICell, R IRow[C]]() IResponse[C, R] {
	return &Response[C, R]{
		RequestTimes: &requestTimes{},
		DataAge:      &dataAge{},
		StatusCode:   http.StatusOK,
		Errors:       []error{},
		Data:         NewTable[C, R](),
		Metadata:     map[string]interface{}{},
	}
}

// ToJson converts a response into json friendly []bye that is indented for readability.
// This is used for passing the data back from the api
func ToJson[C ICell, R IRow[C]](r IResponse[C, R]) (content []byte, err error) {
	return json.MarshalIndent(r, "", "  ")
}

// FromJson converts a []byte back into an IResponse by using json unmarshaling
func FromJson[C ICell, R IRow[C]](content []byte, r IResponse[C, R]) (err error) {
	err = json.Unmarshal(content, r)
	return
}

// FromHttp is similar to FromJson, but first fetches the content from the http.Repsonse body
// and then converts using that
func FromHttp[C ICell, R IRow[C]](content *http.Response, r IResponse[C, R]) (err error) {
	_, bytes := Stringify(content)
	return FromJson[C, R](bytes, r)
}

// Stringify returns the body content of a http.Response as both a string and []byte.
// Very helpful for debugging, testing and converting back and forth from the api.
func Stringify(r *http.Response) (s string, b []byte) {
	b, _ = io.ReadAll(r.Body)
	s = string(b)
	return
}

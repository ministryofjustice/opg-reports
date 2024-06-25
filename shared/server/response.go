package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"opg-reports/shared/data"
	"strings"
	"time"
)

const ResponseTypeHeader string = "X-API-RES-TYPE"

// IApiResponseTimes handles simple start, end and duration elements of the interface.
type IApiResponseTimes interface {
	Start()
	End()
}

// IApiResponseStatus handles tracking the http status of the api response.
// Its value should be used with IApi.Write call at the end
type IApiResponseStatus interface {
	SetStatus(status int)
	GetStatus() int
}

// IApiResponseErrors allows tracking of server side errors such as validation
// and will be included in the IApi.Write
type IApiResponseErrors interface {
	SetErrors(errors []error)
	AddError(err error)
	GetErrors() []error
}

// IApiResponseBase is a merge interface that wuld be typical of an http response.
// This version excludes any result data / handling for simplicty on errors or
// empty results
type IApiResponseBase interface {
	IApiResponseTimes
	IApiResponseStatus
	IApiResponseErrors
	AddErrorWithStatus(err error, status int)
}

// IApiResponseConstraint is used as a constraint on IApiResponseResult to determine
// what form the response data is in, allowing different endpoints to send data back
// is a mix of types.
// For now, we support map[string]R, map[string][]R, map[string]map[string][]R and []R so the "result" field will
// always be either a map or a slice
// Used to differ between a data sequence that is grouped by a key (such as costs) versus
// uptime data which is just a list
type IApiResponseConstraint[R data.IEntry] interface {
	map[string]R | map[string][]R | map[string]map[string][]R | []R
}

// IApiResponseResult providers a response interface whose result type can vary between
// slice, a map or a map of slices.
// This allows api respsones to adapt to the most useful data type for the endpoint
type IApiResponseResult[R data.IEntry, C IApiResponseConstraint[R]] interface {
	IApiResponseBase
	SetResult(result C)
	GetResult() C
	SetType()
}

// ApiResponseConstraintString takes a string version of a contract
// IApiResponseResult constraints for an ApiResponse such as
// map[string][]*cost.Cost and will return them in parts
type ApiResponseConstraintString string

func (r ApiResponseConstraintString) Parts() []string {
	str := string(r)
	return strings.Split(str, "*")
}
func (r ApiResponseConstraintString) C() (g string) {
	parts := r.Parts()
	if len(parts) > 0 {
		g = parts[0]
	}
	return
}
func (r ApiResponseConstraintString) R() (g string) {
	parts := r.Parts()
	if len(parts) > 1 {
		g = parts[1]
	}
	return
}

// ApiResponseTimings impliments [IApiResponseTimes]
type ApiResponseTimings struct {
	Times struct {
		Start    time.Time     `json:"start"`
		End      time.Time     `json:"end"`
		Duration time.Duration `json:"duration"`
	} `json:"timings"`
}

// Start tracks the start time of this request
func (i *ApiResponseTimings) Start() {
	i.Times.Start = time.Now().UTC()
}

// End tracks the end time and the duration of the request
func (i *ApiResponseTimings) End() {
	i.Times.End = time.Now().UTC()
	i.Times.Duration = i.Times.End.Sub(i.Times.Start)
}

// ApiResponseStatus impliments [IApiResponseStatus]
// Provides http status tracking
type ApiResponseStatus struct {
	Status int `json:"status"`
}

// SetStatus updates the status field
func (i *ApiResponseStatus) SetStatus(status int) {
	i.Status = status
}

// GetStatus returns the status field
func (i *ApiResponseStatus) GetStatus() int {
	return i.Status
}

// ApiResponseErrors handles error tracking and impliments [IApiResponseErrors]
type ApiResponseErrors struct {
	Errors []error `json:"errors"`
}

// SetErrors replaces errors with those passed
func (r *ApiResponseErrors) SetErrors(errors []error) {
	r.Errors = errors
}

// AddError add a new error to the list
func (r *ApiResponseErrors) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

// GetErrors returns all errors
func (r *ApiResponseErrors) GetErrors() []error {
	return r.Errors
}

// ApiResponseBase impliments IApiResponseBase
// Would be used for a simple endpoint that doesn't return data,
// such as an api root
type ApiResponseBase struct {
	ApiResponseTimings
	ApiResponseStatus
	ApiResponseErrors
}

// AddErrorWithStatus adds an error and updates the status at the same time.
// Helpful when validating fields to do both at once.
func (i *ApiResponseBase) AddErrorWithStatus(err error, status int) {
	i.AddError(err)
	i.SetStatus(status)
}

// NewSimpleApiResponse returns a fresh ApiResponseBase with
// status set as OK and errors empty
func NewSimpleApiResponse() *ApiResponseBase {
	return &ApiResponseBase{
		ApiResponseTimings: ApiResponseTimings{},
		ApiResponseStatus:  ApiResponseStatus{Status: http.StatusOK},
		ApiResponseErrors:  ApiResponseErrors{Errors: []error{}},
	}
}

// ApiResponse impliments [IApiResponse].
// It allows a response to return with variable (C) data type. This is currently
// constrained to map[string]R, map[string][]R and []R.
// This means various enpoints can return differing ways collecting the data.IEntry
// so some can group by a field or just list everything that matches
//
// This struct and interface allows you to easily decode a response as long as you know
// its return type
type ApiResponse[R data.IEntry, C IApiResponseConstraint[R]] struct {
	ApiResponseBase
	Type   string `json:"result_type"`
	Result C      `json:"result"`
}

// SetResult updates the internal result data
func (i *ApiResponse[R, C]) SetResult(result C) {
	i.Result = result
}

// GetResult returns the result
func (i *ApiResponse[R, C]) GetResult() C {
	return i.Result
}

// SetType updates the Type field to a name for easier tracking / decoding
// for the response
func (i *ApiResponse[R, C]) SetType() {
	var x C
	i.Type = fmt.Sprintf("%T", x)
}

// NewApiResponse returns an ApiResponse that handles results of the types set
func NewApiResponse[R data.IEntry, C IApiResponseConstraint[R]]() *ApiResponse[R, C] {
	return &ApiResponse[R, C]{
		ApiResponseBase: *NewSimpleApiResponse(),
	}
}

// NewApiResponseFromJson unmarshals the content into a response of types set
func NewApiResponseFromJson[R data.IEntry, C IApiResponseConstraint[R]](content []byte) (response *ApiResponse[R, C], err error) {
	err = json.Unmarshal(content, &response)
	return
}

// ResponseAsStrings takes a http.Response and returns string & []byte
// values of the response body
func ResponseAsStrings(r *http.Response) (string, []byte) {
	b, _ := io.ReadAll(r.Body)
	return string(b), b
}

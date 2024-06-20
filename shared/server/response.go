package server

import (
	"encoding/json"
	"net/http"
	"opg-reports/shared/data"
	"time"
)

type IApiResponseTimes interface {
	Start()
	End()
}

type IApiResponseStatus interface {
	SetStatus(status int)
	GetStatus() int
}

type IApiResponseErrors interface {
	SetErrors(errors []error)
	AddError(err error)
	GetErrors() []error
}

type IApiResponseBody interface {
	Body() []byte
}

type IApiResponseBase interface {
	IApiResponseTimes
	IApiResponseStatus
	IApiResponseErrors
	IApiResponseBody
	AddErrorWithStatus(err error, status int)
}

type IApiResponseResultConstraint[R data.IEntry] interface {
	map[string]R | map[string][]R | []R
}

type IApiResponseResult[R data.IEntry, C IApiResponseResultConstraint[R]] interface {
	IApiResponseBase
	SetResult(result C)
	GetResult() C
}

// ApiResponseTimings impliments [IApiResponseTimes]
type ApiResponseTimings struct {
	Times struct {
		Start    time.Time     `json:"start"`
		End      time.Time     `json:"end"`
		Duration time.Duration `json:"duration"`
	} `json:"timings"`
}

func (i *ApiResponseTimings) Start() {
	i.Times.Start = time.Now().UTC()
}
func (i *ApiResponseTimings) End() {
	i.Times.End = time.Now().UTC()
	i.Times.Duration = i.Times.End.Sub(i.Times.Start)
}

// ApiResponseStatus
type ApiResponseStatus struct {
	Status int `json:"status"`
}

func (i *ApiResponseStatus) SetStatus(status int) {
	i.Status = status
}
func (i *ApiResponseStatus) GetStatus() int {
	return i.Status
}

// ApiResponseErrors
type ApiResponseErrors struct {
	Errors []error `json:"errors"`
}

func (r *ApiResponseErrors) SetErrors(errors []error) {
	r.Errors = errors
}

func (r *ApiResponseErrors) AddError(err error) {
	r.Errors = append(r.Errors, err)
}
func (r *ApiResponseErrors) GetErrors() []error {
	return r.Errors
}

// ApiResponseBody
type ApiResponseBody struct{}

func (r *ApiResponseBody) Body() []byte {
	body, _ := json.Marshal(r)
	return body
}

// IApiResponseBase
type ApiResponseBase struct {
	ApiResponseTimings
	ApiResponseStatus
	ApiResponseErrors
	ApiResponseBody
}

func (i *ApiResponseBase) AddErrorWithStatus(err error, status int) {
	i.AddError(err)
	i.SetStatus(status)
}

func NewApiSimpleResponse() *ApiResponseBase {
	return &ApiResponseBase{
		ApiResponseTimings: ApiResponseTimings{},
		ApiResponseStatus:  ApiResponseStatus{Status: http.StatusOK},
		ApiResponseErrors:  ApiResponseErrors{Errors: []error{}},
		ApiResponseBody:    ApiResponseBody{},
	}
}

type ApiResponseWithResults[R data.IEntry, C IApiResponseResultConstraint[R]] struct {
	ApiResponseBase
	Result C `json:"result"`
}

func (i *ApiResponseWithResults[R, C]) SetResult(result C) {
	i.Result = result
}

func (i *ApiResponseWithResults[R, C]) GetResult() C {
	return i.Result
}

func NewApiResponseWithResult[R data.IEntry, C IApiResponseResultConstraint[R]]() *ApiResponseWithResults[R, C] {
	return &ApiResponseWithResults[R, C]{
		ApiResponseBase: *NewApiSimpleResponse(),
	}
}

// type ApiResponseResultIsMapSlice[R data.IEntry, C map[string][]R] struct {
// 	ApiResponseBase
// 	Result C `json:"result"`
// }

// func (i *ApiResponseResultIsMapSlice[R, C]) SetResult(result C) {
// 	i.Result = result
// }

// func (i *ApiResponseResultIsMapSlice[R, C]) GetResult() C {
// 	return i.Result
// }

// type ApiResponseResultIsSlice[R data.IEntry, C []R] struct {
// 	ApiResponseBase
// 	Result C `json:"result"`
// }

// func (i *ApiResponseResultIsSlice[R, C]) SetResult(result C) {
// 	i.Result = result
// }

// func (i *ApiResponseResultIsSlice[R, C]) GetResult() C {
// 	return i.Result
// }

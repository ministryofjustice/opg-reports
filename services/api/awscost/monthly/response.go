package monthly

import (
	"encoding/json"
	"time"
)

type ApiResponse struct {
	RequestStart    time.Time     `json:"request_start"`
	RequestEnd      time.Time     `json:"request_end"`
	RequestDuration time.Duration `json:"request_duration"`
	Errors          []error       `json:"errors"`
	Status          int           `json:"status"`
	Result          interface{}   `json:"result"`
}

func (r *ApiResponse) Body() []byte {
	body, _ := json.Marshal(r)
	return body
}

func (r *ApiResponse) Start() {
	r.RequestStart = time.Now().UTC()
}

func (r *ApiResponse) End() {
	r.RequestEnd = time.Now().UTC()
	r.RequestDuration = r.RequestEnd.Sub(r.RequestStart)
}

func (r *ApiResponse) SetResults(results interface{}) {
	r.Result = results
}

func (r *ApiResponse) GetResults() interface{} {
	return r.Result
}

func (r *ApiResponse) SetStatus(status int) {
	r.Status = status
}
func (r *ApiResponse) GetStatus() int {
	return r.Status
}

func (r *ApiResponse) SetErrors(errors []error) {
	r.Errors = errors
}

func (r *ApiResponse) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

func (r *ApiResponse) AddStatusError(status int, err error) {
	r.Errors = append(r.Errors, err)
	r.SetStatus(status)
}

func (r *ApiResponse) GetErrors() []error {
	return r.Errors
}

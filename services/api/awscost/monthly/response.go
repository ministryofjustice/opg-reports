package monthly

import (
	"encoding/json"
	"time"
)

type ApiResponse struct {
	RequestStart    time.Time     `json:"request_start"`
	RequestEnd      time.Time     `json:"request_end"`
	RequestDuration time.Duration `json:"request_duration"`
	Errs            []error       `json:"errors"`
	StatusCode      int           `json:"status"`
	Res             interface{}   `json:"result"`
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
func (r *ApiResponse) Set(results interface{}, status int, errors []error) {
	r.Res = results
	r.StatusCode = status
	r.Errs = errors
}
func (r *ApiResponse) Results() interface{} {
	return r.Res
}

func (r *ApiResponse) Status() int {
	return r.StatusCode
}

func (r *ApiResponse) Errors() []error {
	return r.Errs
}

package resp

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type RequestTimings struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration float64   `json:"duration"`
}

type DataAge struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type Response struct {
	Timer      *RequestTimings          `json:"request_timings,omitempty"`
	DataAge    *DataAge                 `json:"data_age"`
	StatusCode int                      `json:"status"`
	Errors     []error                  `json:"errors"`
	Metadata   map[string]interface{}   `json:"metadata"`
	Result     []map[string]interface{} `json:"result"`
}

func (rp *Response) Start(w http.ResponseWriter, r *http.Request) {
	rp.StatusCode = http.StatusOK

	slog.Info("request start",
		slog.String("request_method", r.Method),
		slog.String("request_uri", r.URL.String()))
	rp.TimerStart()
}

func (rp *Response) End(w http.ResponseWriter, r *http.Request) {
	rp.TimerEnd()
	rp.TimerDuration()

	content, err := json.MarshalIndent(rp, "", "  ")
	if err != nil {
		rp.Errors = append(rp.Errors, err)
		slog.Error(err.Error())
	}

	// set default error header
	if len(rp.Errors) > 0 && rp.StatusCode == http.StatusOK {
		rp.StatusCode = http.StatusBadRequest
	}

	slog.Info("request end",
		slog.Int("status", rp.StatusCode),
		slog.Int("errors", len(rp.Errors)),
		slog.String("request_method", r.Method),
		slog.String("request_uri", r.URL.String()))

	w.WriteHeader(rp.StatusCode)
	w.Write(content)
}

// TimerStart is called early in the http request handler to
// accurately track the start time of the request.
// This is then used with end time to work out durations
// that can be checked for performance etc
func (r *Response) TimerStart() {
	r.Timer.Start = time.Now().UTC()
}

// TimerEnd is companiion to TimerStart, this tracks the
// end of the http request being processed and can then
// be used to work out duration.
func (r *Response) TimerEnd() {
	r.Timer.End = time.Now().UTC()
}

// TimerDuration uses the end & start times to work out how long
// a http request has taken to be processed. This informaton
// is helpful for assessing performance over the api
func (r *Response) TimerDuration() time.Duration {
	dur := r.Timer.End.Sub(r.Timer.Start)
	r.Timer.Duration = dur.Seconds()
	return dur
}

func New() *Response {
	return &Response{
		Timer:      &RequestTimings{},
		DataAge:    &DataAge{},
		StatusCode: http.StatusOK,
		Errors:     []error{},
		Result:     []map[string]interface{}{},
		Metadata:   map[string]interface{}{},
	}
}

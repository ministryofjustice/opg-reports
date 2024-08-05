package resp

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"opg-reports/shared/server/resp/table"
	"slices"
	"time"
)

type RequestTimings struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

type DataAge struct {
	Min *time.Time `json:"min"`
	Max *time.Time `json:"max"`
	All []int64    `json:"-"`
}

type Response struct {
	Timer      *RequestTimings        `json:"request_timings,omitempty"`
	DataAge    *DataAge               `json:"data_age"`
	StatusCode int                    `json:"status"`
	Errors     []error                `json:"errors"`
	Metadata   map[string]interface{} `json:"metadata"`
	Result     *table.Table           `json:"result"`
}

func (rp *Response) Start(w http.ResponseWriter, r *http.Request) {
	slog.Info("request start",
		slog.String("request_method", r.Method),
		slog.String("request_uri", r.URL.String()))
	rp.TimerStart()
}

func (rp *Response) End(w http.ResponseWriter, r *http.Request) {
	rp.TimerEnd()
	rp.TimerDuration()
	rp.GetDataAgeMin()
	rp.GetDataAgeMax()

	content, _ := json.MarshalIndent(rp, "", "  ")

	slog.Info("request end",
		slog.Int("status", rp.StatusCode),
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
	r.Timer.Duration = r.Timer.End.Sub(r.Timer.Start)
	return r.Timer.Duration
}

// AddDataAge is used to track the created times of the data
// associated with this api response. This provides a way
// to show messages about when the data was created / updated
// in front ends
func (r *Response) AddDataAge(times ...time.Time) {
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
func (r *Response) GetDataAgeMin() (t time.Time) {
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
func (r *Response) GetDataAgeMax() (t time.Time) {
	if r.DataAge.Max != nil {
		return *r.DataAge.Max
	} else if len(r.DataAge.All) > 0 {
		max := slices.Max(r.DataAge.All)
		t = time.UnixMicro(max).UTC()
		r.DataAge.Max = &t
	}
	return t
}

func New() *Response {
	return &Response{
		Timer:      &RequestTimings{},
		DataAge:    &DataAge{},
		StatusCode: http.StatusOK,
		Errors:     []error{},
		Metadata:   map[string]interface{}{},
		Result:     table.New(),
	}
}

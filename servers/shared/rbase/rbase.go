package rbase

import (
	"context"
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

type Respond interface {
	StatusGet() int
	StatusSet(status int)
	ErrorAdd(err error)
	ErrorGet() []string
	Timer() *RequestTimings
}

type RequestTimings struct {
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Seconds float64   `json:"duration"`
}

func (rt *RequestTimings) Duration() time.Duration {
	dur := rt.End.Sub(rt.Start)
	rt.Seconds = dur.Seconds()
	return dur
}

type DataAge struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type Response struct {
	ctx          context.Context
	RequestTimer *RequestTimings `json:"request_timings,omitempty"`
	DataAge      *DataAge        `json:"data_age"`
	StatusCode   int             `json:"status"`
	Errors       []string        `json:"errors"`
	StartDate    string          `json:"start_date"`
	EndDate      string          `json:"end_date"`
	DateRange    []string        `json:"date_range"`

	// hide all of these, only used on the front reader
	RowFilters       map[string]interface{}                `json:"row_filters,omitempty"`
	Organisation     string                                `json:"-"`
	PageTitle        string                                `json:"-"`
	NavigationTop    map[string]*navigation.NavigationItem `json:"-"`
	NavigationSide   []*navigation.NavigationItem          `json:"-"`
	NavigationActive *navigation.NavigationItem            `json:"-"`
	Rows             map[string]map[string]interface{}     `json:"-"`
}

func (rp *Response) StatusGet() int {
	return rp.StatusCode
}
func (rp *Response) StatusSet(status int) {
	rp.StatusCode = status
}
func (rp *Response) Timer() *RequestTimings {
	return rp.RequestTimer
}
func (rp *Response) ErrorAdd(err error) {
	rp.Errors = append(rp.Errors, err.Error())
}
func (rp *Response) ErrorGet() []string {
	return rp.Errors
}

// ----

func Start[T Respond](rp T, w http.ResponseWriter, r *http.Request) {
	rp.StatusSet(http.StatusOK)
	ts := rp.Timer()
	ts.Start = time.Now().UTC()
}

func Stop[T Respond](rp T, w http.ResponseWriter, r *http.Request) {
	rp.Timer().End = time.Now().UTC()
	rp.Timer().Duration()
}

func End[T Respond](rp T, w http.ResponseWriter, r *http.Request) {
	Stop(rp, w, r)
	content, err := convert.Marshal(rp)
	if err != nil {
		rp.ErrorAdd(err)
	}

	if len(rp.ErrorGet()) > 0 && rp.StatusGet() == http.StatusOK {
		rp.StatusSet(http.StatusBadRequest)
	}
	w.WriteHeader(rp.StatusGet())
	w.Write(content)
}

func ErrorAndEnd[T Respond](rp T, err error, w http.ResponseWriter, r *http.Request) {
	rp.ErrorAdd(err)
	End(rp, w, r)
}

package blocks

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/api"
	"opg-reports/report/internal/domain/codeowners/codeownerapis/codeownerdynamic"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/domain/infracosts/infracostapis/infracostdynamic"
	"opg-reports/report/internal/domain/teams/teamapis/teamdynamic"
	"opg-reports/report/internal/domain/uptime/uptimeapis/uptimedynamic"
	"opg-reports/report/internal/utils/times"
	"time"
)

var timeout = (2 * time.Second)

// TeamData calls the api and fetches team data
func TeamData(ctx context.Context, log *slog.Logger, apiHost string, request *http.Request) (teams []string, err error) {
	var (
		resp *teamdynamic.TeamResponseBody
		ep   string = "/v1/teams"
	)
	teams = []string{}
	resp, _, err = api.Get[*teamdynamic.TeamResponseBody](ctx, log, &api.Call{
		Host:     apiHost,
		Endpoint: ep,
		Request:  request,
		Timeout:  timeout,
	})
	if err != nil {
		return
	}

	for _, team := range resp.Data {
		teams = append(teams, team.Name)
	}

	return
}

// UptimeData gets uptime info from the api
func UptimeData(ctx context.Context, log *slog.Logger, apiHost string, request *http.Request, overwrites ...*api.Param) (data []map[string]interface{}, headers map[string][]string, err error) {
	var (
		resp   *uptimedynamic.UptimeResponseBody
		params []*api.Param = []*api.Param{}
		ep     string       = `/v1/uptime/between/{start_date}/{end_date}`
		now    time.Time    = times.Today()
		start  time.Time    = times.Ago(times.ResetMonth(now), 6, times.MONTH)
	)
	data = []map[string]interface{}{}
	headers = map[string][]string{}

	params = []*api.Param{
		{Type: api.PATH, Key: "start_date", Value: times.AsYMString(start)},
		{Type: api.PATH, Key: "end_date", Value: times.AsYMString(now)},
		{Type: api.QUERY, Key: "team", Value: ""},
	}
	// overwrite with optionals
	params = overwriteParams(params, overwrites...)

	resp, _, err = api.Get[*uptimedynamic.UptimeResponseBody](ctx, log, &api.Call{
		Host:     apiHost,
		Endpoint: ep,
		Request:  request,
		Params:   params,
		Timeout:  timeout,
	})
	if err != nil {
		return
	}
	data = resp.Data
	headers = resp.Headers

	return
}

// InfracostData gets cost info from the api
func InfracostData(ctx context.Context, log *slog.Logger, apiHost string, request *http.Request, overwrites ...*api.Param) (data []map[string]interface{}, headers map[string][]string, err error) {
	var (
		resp   *infracostdynamic.InfracostResponseBody
		params []*api.Param = []*api.Param{}
		ep     string       = `/v1/infacosts/between/{start_date}/{end_date}`
		now    time.Time    = times.Today()
		start  time.Time    = times.Ago(times.ResetMonth(now), 5, times.MONTH)
	)
	data = []map[string]interface{}{}
	headers = map[string][]string{}
	// the base line params
	params = []*api.Param{
		{Type: api.PATH, Key: "start_date", Value: times.AsYMString(start)},
		{Type: api.PATH, Key: "end_date", Value: times.AsYMString(now)},
		{Type: api.QUERY, Key: "team", Value: ""},
		{Type: api.QUERY, Key: "account", Value: ""},
		{Type: api.QUERY, Key: "environment", Value: ""},
		{Type: api.QUERY, Key: "service", Value: ""},
		{Type: api.QUERY, Key: "sort", Value: ""},
	}
	// overwrite with optionals
	params = overwriteParams(params, overwrites...)

	resp, _, err = api.Get[*infracostdynamic.InfracostResponseBody](ctx, log, &api.Call{
		Host:     apiHost,
		Endpoint: ep,
		Request:  request,
		Params:   params,
		Timeout:  timeout,
	})
	if err != nil {
		return
	}
	data = resp.Data
	headers = resp.Headers

	return
}

// CodeownerData gets codeowners without any owners from the api
func CodeownerData(ctx context.Context, log *slog.Logger, apiHost string, request *http.Request, overwrites ...*api.Param) (data []*codeownermodels.CodeownerData, err error) {
	var (
		resp   *codeownerdynamic.CodeownerResponseBody
		params []*api.Param = []*api.Param{}
		ep     string       = `/v1/codeowners`
	)
	data = []*codeownermodels.CodeownerData{}
	params = []*api.Param{
		{Type: api.QUERY, Key: "team", Value: ""},
		{Type: api.QUERY, Key: "account", Value: ""},
		{Type: api.QUERY, Key: "codeowner", Value: ""},
		{Type: api.QUERY, Key: "codebase", Value: ""},
	}
	// overwrite with optionals
	params = overwriteParams(params, overwrites...)

	resp, _, err = api.Get[*codeownerdynamic.CodeownerResponseBody](ctx, log, &api.Call{
		Host:     apiHost,
		Endpoint: ep,
		Request:  request,
		Params:   params,
		Timeout:  timeout,
	})
	if err != nil {
		return
	}
	data = resp.Data

	return
}

// overwriteParams makes sure the extras are added to the main list and any default params
// that match are skipped - allows replace teams with a locked version etc
func overwriteParams(params []*api.Param, overwrites ...*api.Param) (list []*api.Param) {
	list = []*api.Param{}
	for _, ex := range overwrites {
		list = append(list, ex)
	}
	for _, p := range params {
		var added = false
		for _, ex := range list {
			if ex.Key == p.Key {
				added = true
			}
		}
		if !added {
			list = append(list, p)
		}
	}
	return
}

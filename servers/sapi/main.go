/*
sapi runs the api to surface data from the registered endpoints and way to output the spec.

Usage:

	sapi
	spai openapi

Calling `openapi` will output to stdout the yaml spec for the api.

To expand this api with new content, please look at how `costsapi` is setup and when you
have an equilivant append this to the `segments` map.

Registered segments:
  - costsapi
*/
package main

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/ministryofjustice/opg-reports/pkg/bi"
	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/envar"
	"github.com/ministryofjustice/opg-reports/servers/sapi/lib"
	"github.com/ministryofjustice/opg-reports/sources/costs"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsapi"
	"github.com/ministryofjustice/opg-reports/sources/standards"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsapi"
	"github.com/ministryofjustice/opg-reports/sources/uptime"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimeapi"
)

var mode = bi.Mode

// we split the api handlers into simple & full groups
// `simple` is used for the basic install
// `full` covers all options
// Set using the bi.Mode which is a ldflag
var (
	simpleSegments map[string]*lib.ApiSegment = map[string]*lib.ApiSegment{
		standardsapi.Segment: {
			DbFile:       "./databases/standards.db",
			SetupFunc:    standards.Setup,
			RegisterFunc: standardsapi.Register,
		},
	}
	fullSegments map[string]*lib.ApiSegment = map[string]*lib.ApiSegment{
		costsapi.Segment: {
			DbFile:       "./databases/costs.db",
			SetupFunc:    costs.Setup,
			RegisterFunc: costsapi.Register,
		},
		standardsapi.Segment: {
			DbFile:       "./databases/standards.db",
			SetupFunc:    standards.Setup,
			RegisterFunc: standardsapi.Register,
		},
		uptimeapi.Segment: {
			DbFile:       "./databases/uptime.db",
			SetupFunc:    uptime.Setup,
			RegisterFunc: uptimeapi.Register,
		},
	}
	segmentChoices map[string]map[string]*lib.ApiSegment = map[string]map[string]*lib.ApiSegment{
		"simple": simpleSegments,
		"full":   fullSegments,
	}
	segments = segmentChoices[mode]
)

// init is used to fetch stored databases from s3
// or create dummy versions of them
func init() {
	var ctx context.Context = context.Background()
	lib.SetupSegments(ctx, segments)
}

// main executes the clis wrapped huma api
func main() {
	bi.Dump()
	Run()
}

// Run is the main execution loop
// It gets the cli from inside lib
func Run() {
	var (
		api        huma.API
		server     http.Server
		mux        *http.ServeMux  = http.NewServeMux()
		ctx        context.Context = context.Background()
		apiTitle   string          = lib.ApiTitle()
		apiVersion string          = lib.ApiVersion()
		addr       string          = envar.Get("API_ADDR", consts.ServerDefaultApiAddr)
	)
	// create the server
	server = http.Server{
		Addr:    addr,
		Handler: mux,
	}
	// create the api
	api = humago.New(mux, huma.DefaultConfig(apiTitle, apiVersion))
	// get the cli and run it
	cmd := lib.CLI(ctx, api, &server, segments)
	cmd.Run()

}

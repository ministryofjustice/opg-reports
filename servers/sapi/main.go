/*
sapi runs the api to surface data from the registered endpoints and way to output the spec.

Usage:

	sapi [flags]
	spai openapi

The flags are:

	-port=<value>
		The port the server will start running on.
		Default: 8081
	-host=<value>
		The host the server is bound to.
		Default: localhost
	-debug
		Enable more detailed logging.
		CURRENTLY NOT FUNCTIONAL

Calling `openapi` will output to stdout the yaml spec for the api.

To expand this api with new content, please look at how `costsapi` is setup and when you
have an equilivant append this to the `segments` map.

Registered segments:
  - costsapi
*/
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/ministryofjustice/opg-reports/pkg/bi"
	"github.com/ministryofjustice/opg-reports/pkg/envar"
	"github.com/ministryofjustice/opg-reports/servers/sapi/lib"
	"github.com/ministryofjustice/opg-reports/sources/costs"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsapi"
)

var segments map[string]*lib.ApiSegment = map[string]*lib.ApiSegment{

	costsapi.Segment: {
		DbFile:       "./databases/costs.db",
		SetupFunc:    costs.Setup,
		RegisterFunc: costsapi.Register,
	},
}

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
		port       string          = envar.Get("API_PORT", "8081")
	)
	// create the server
	server = http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: mux,
	}
	// create the api
	api = humago.New(mux, huma.DefaultConfig(apiTitle, apiVersion))
	// get the cli and run it
	cmd := lib.CLI(ctx, api, &server, segments)
	cmd.Run()

}

/*
sapi runs the api to surface data from the registered endpoints and way to output the spec.

Usage:

	api
	pai openapi

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
	"github.com/ministryofjustice/opg-reports/info"
	"github.com/ministryofjustice/opg-reports/internal/envar"
	"github.com/ministryofjustice/opg-reports/internal/fileutils"
	"github.com/ministryofjustice/opg-reports/servers/api/handlers"
	"github.com/ministryofjustice/opg-reports/servers/api/lib"
)

var (
	mode        string = info.Fixtures        // decides which set of endpoints to use
	localDBPath string = "./databases/api.db" // path to the database
	bucketName  string = info.BucketName      // name of the bucket to pull database from
	bucketDB    string = "./databases/api.db" // name of the database within the bucket
)

var (
	allHandlers map[string]map[string]lib.RegisterHandlerFunc = map[string]map[string]lib.RegisterHandlerFunc{
		"simple": {
			"dataset":                     handlers.RegisterDatasets,
			"github-repository-standards": handlers.RegisterGitHubRepositoryStandards,
			"github-teams":                handlers.RegisterGitHubTeams,
			"units":                       handlers.RegisterUnits,
		},
		"full": {
			"aws-accounts":                handlers.RegisterAwsAccounts,
			"aws-costs":                   handlers.RegisterAwsCosts,
			"aws-uptime":                  handlers.RegisterAwsUptime,
			"dataset":                     handlers.RegisterDatasets,
			"github-repositories":         handlers.RegisterGitHubRepositories,
			"github-releases":             handlers.RegisterGitHubRelases,
			"github-repository-standards": handlers.RegisterGitHubRepositoryStandards,
			"github-teams":                handlers.RegisterGitHubTeams,
			"units":                       handlers.RegisterUnits,
		},
	}
	Handlers map[string]lib.RegisterHandlerFunc = allHandlers[mode]
)

// init is used to fetch stored databases from s3
// or create dummy versions of them
func init() {
	var ctx context.Context = context.Background()

	// new way of seeding
	// - if we are using the real data set, go fetch it
	if info.Dataset == "real" {
		lib.DownloadS3DB(bucketName, bucketDB, localDBPath)
	}
	// if the local db does not exist, then create a seeded version
	if !fileutils.Exists(localDBPath) {
		lib.SeedDB(ctx, localDBPath)
	}

}

// main executes the clis wrapped huma api
func main() {
	info.Log()
	Run()
}

// Run is the main execution loop
// It gets the cli from inside lib
func Run() {
	var (
		api          huma.API
		server       http.Server
		mux          *http.ServeMux  = http.NewServeMux()
		ctx          context.Context = context.Background()
		apiTitle     string          = lib.ApiTitle()
		apiVersion   string          = lib.ApiVersion()
		addr         string          = envar.Get("API_ADDR", info.ServerDefaultApiAddr)
		databasePath string          = envar.Get("DB_PATH", localDBPath)
	)

	// create the server
	server = http.Server{
		Addr:    addr,
		Handler: mux,
	}
	// create the api
	api = humago.New(mux, huma.DefaultConfig(apiTitle, apiVersion))
	// get the cli and run it
	cmd := lib.CLI(ctx, api, &server, Handlers, databasePath)
	cmd.Run()

}

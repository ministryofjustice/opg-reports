package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/ministryofjustice/opg-reports/api/awscosts"
	"github.com/ministryofjustice/opg-reports/versions"
)

type HomepageResponse struct {
	Body struct {
		Message string `json:"message" example:"Successful connection."`
	}
}

// Opts provides a series of values for the api server that is configured
// at run time
type Opts struct {
	Debug bool   `doc:"When true enables more detailed logging." default:"false"`
	Host  string `doc:"Hostname to listen on." default:"localhost"`
	Port  int    `doc:"Port to listen on." default:"8081"`
	Spec  bool   `doc:"When true, the openapi spec will be show on server startup" default:"false"`
}

type apiRegistrationFunc func(api huma.API, dbPath string)
type apiSetupFunc func(ctx context.Context) (err error)

type apiSegment struct {
	DbFile       string
	RegisterFunc apiRegistrationFunc
	SetupFunc    apiSetupFunc
}

var segments map[string]*apiSegment = map[string]*apiSegment{
	"awscosts": {DbFile: "./dbs/awscosts.db", RegisterFunc: awscosts.Register, SetupFunc: awscosts.Setup},
}

func main() {

	cli := humacli.New(func(hooks humacli.Hooks, opts *Opts) {

		// create the server
		mux := http.NewServeMux()
		server := http.Server{
			Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
			Handler: mux,
		}

		// create the api
		versionStr := fmt.Sprintf("%s [%s] (%s)", versions.Build, versions.Timestamp, versions.Commit)
		api := humago.New(mux, huma.DefaultConfig("Reporting API", versionStr))

		// register homepage action that will return an almost empty result
		huma.Register(api, huma.Operation{
			OperationID:   "get-homepage",
			Method:        http.MethodGet,
			Path:          "/",
			Summary:       "Home",
			Description:   "Operates as the root of the API and a simple endpoint to test connectivity with, but returns no reporting data.",
			DefaultStatus: http.StatusOK,
		}, func(ctx context.Context, input *struct{}) (homepage *HomepageResponse, err error) {
			homepage = &HomepageResponse{}
			homepage.Body.Message = "Successful connection"
			return
		})

		if !opts.Spec {
			// Register the sub helpers
			for name, segment := range segments {
				slog.Info("[api.main] registering", slog.String("segment", name))
				segment.RegisterFunc(api, segment.DbFile)
			}
		}
		// run the server or show the spec
		hooks.OnStart(func() {
			if opts.Spec {
				bytes, _ := api.OpenAPI().YAML()
				fmt.Println(string(bytes))
			} else {
				// output info that the server is starting
				slog.Info("Starting api server.",
					slog.Bool("debug", opts.Debug),
					slog.Bool("spec", opts.Spec),
					slog.String("host", opts.Host),
					slog.Int("port", opts.Port))
				slog.Info(fmt.Sprintf("Docs: [http://%s:%d/docs]", opts.Host, opts.Port))
				server.ListenAndServe()
			}
		})

		// graceful shutdown
		hooks.OnStop(func() {
			var shutdownDelay = 3 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), shutdownDelay)
			defer cancel()
			server.Shutdown(ctx)
		})

	})

	cli.Run()
}

// init is used to fetch stored databases from s3
// or create dummy versions of them
func init() {
	var ctx context.Context = context.Background()
	var err error

	for name, segment := range segments {
		slog.Debug("[api.init] running setup", slog.String("segment", name))

		if err = segment.SetupFunc(ctx); err != nil {
			panic(err)
		}

	}
}

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
	"github.com/ministryofjustice/opg-reports/costs"
	"github.com/ministryofjustice/opg-reports/costs/costsapi"
	"github.com/ministryofjustice/opg-reports/pkg/bi"
	"github.com/spf13/cobra"
)

var segments map[string]*apiSegment = map[string]*apiSegment{

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

	for name, segment := range segments {
		slog.Info("[api.init]", slog.String("segment", name))
		segment.SetupFunc(ctx, segment.DbFile)
	}
}

// main executes the clis wrapped huma api
func main() {
	slog.Info("Build info",
		slog.String("ApiVersion", bi.ApiVersion),
		slog.String("Commit", bi.Commit),
		slog.String("Organisation", bi.Organisation),
		slog.String("Semver", bi.Semver),
		slog.String("Timestamp", bi.Timestamp),
	)
	run()
}

type apiSegment struct {
	DbFile       string
	SetupFunc    func(ctx context.Context, dbFilepath string)
	RegisterFunc func(api huma.API)
}

// HomepageResponse used for simple / page
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
}

func apiMiddleware(ctx huma.Context, next func(huma.Context)) {
	for segment, cfg := range segments {
		ctx = huma.WithValue(ctx, segment, cfg.DbFile)
	}
	next(ctx)
}

func addBaseSetup(api huma.API) {
	// Register the middleware
	api.UseMiddleware(apiMiddleware)
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
}

func run() {
	var api huma.API
	var shutdownDelay = 5 * time.Second
	var ctx = context.Background()
	var versionStr = fmt.Sprintf("%s [%s] (%s)", bi.Semver, bi.Timestamp, bi.Commit)
	var apiTitle = fmt.Sprintf("%s Reports API", bi.Organisation)

	cli := humacli.New(func(hooks humacli.Hooks, opts *Opts) {
		// create the server
		mux := http.NewServeMux()
		server := http.Server{
			Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
			Handler: mux,
		}

		// create the api
		api = humago.New(mux, huma.DefaultConfig(apiTitle, versionStr))

		addBaseSetup(api)

		slog.Info("[api.main] registering...")
		for name, segment := range segments {
			slog.Info("[api.main] register segment", slog.String("segment", name))
			segment.RegisterFunc(api)
		}
		slog.Info("[api.main] registered.")

		hooks.OnStart(func() {
			slog.Info("Starting api server.",
				slog.Bool("debug", opts.Debug),
				slog.String("host", opts.Host),
				slog.Int("port", opts.Port))
			slog.Info(fmt.Sprintf("API: [http://%s:%d/]", opts.Host, opts.Port))
			slog.Info(fmt.Sprintf("Docs: [http://%s:%d/docs]", opts.Host, opts.Port))

			server.ListenAndServe()
		})

		// graceful shutdown
		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(ctx, shutdownDelay)
			defer cancel()
			server.Shutdown(ctx)
		})

	})
	// Add command to dump out yaml
	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			// Use downgrade to return OpenAPI 3.0.3 YAML since oapi-codegen doesn't
			// support OpenAPI 3.1 fully yet. Use `.YAML()` instead for 3.1.
			b, _ := api.OpenAPI().DowngradeYAML()
			fmt.Println(string(b))
		},
	})

	cli.Run()
}

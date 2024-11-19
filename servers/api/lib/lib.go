package lib

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/ministryofjustice/opg-reports/pkg/bi"
	"github.com/spf13/cobra"
)

// ApiSegment captures data for each part of the api (in the folder ./sources/)
type ApiSegment struct {
	DbFile       string
	SetupFunc    func(ctx context.Context, dbFilepath string, seed bool)
	RegisterFunc func(api huma.API)
}

// CliOptions is empty, used just as a place holder to match the func signature
type CliOptions struct{}

type HomepageResponse struct {
	Body struct {
		Message string `json:"message" example:"Successful connection."`
	}
}

// SetupSegments runs the setup functions for all the segments
func SetupSegments(ctx context.Context, segments map[string]*ApiSegment) {
	for name, segment := range segments {
		slog.Info("[api.init]", slog.String("segment", name))
		segment.SetupFunc(ctx, segment.DbFile, true)
	}
}

// AddMiddleware adds the standard middleware information and process
// for each api segment
// Currently - adds database path as a value to the context
func AddMiddleware(api huma.API, segments map[string]*ApiSegment) {
	//
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		for segment, cfg := range segments {
			ctx = huma.WithValue(ctx, segment, cfg.DbFile)
		}
		next(ctx)
	})

}

// AddHomepage registers a simple home page as the default
// root of the api
func AddHomepage(api huma.API, segments map[string]*ApiSegment) {
	//
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

// RegisterSegments calls the registration function for each of the api segments
// allowing them to attach their own routes to the api
func RegisterSegments(api huma.API, segments map[string]*ApiSegment) {
	slog.Info("[api.RegisterSegments] registering ...")
	for name, segment := range segments {
		slog.Info("[api.RegisterSegments] register segment", slog.String("segment", name))
		segment.RegisterFunc(api)
	}
	slog.Info("[api.RegisterSegments] registered.")
}

// ApiTitle generates the title for the api using build details
func ApiTitle() string {
	return fmt.Sprintf("%s Reports API", bi.Organisation)
}

// ApiVersion generates the version string for the api from build details.
// Currently:
//
//	`<semver> [<timestamp>] (<git-sha>)`
func ApiVersion() string {
	return bi.Signature()
}

// CLI returns the api wrapped as a cli command and appends not only the api routes
// but also a command to output the api spec (`openapi`)
func CLI(ctx context.Context, api huma.API, server *http.Server, segments map[string]*ApiSegment) (cli humacli.CLI) {
	var shutdownDelay time.Duration = 5 * time.Second

	cli = humacli.New(func(hooks humacli.Hooks, opts *CliOptions) {
		var addr = server.Addr

		AddMiddleware(api, segments)
		AddHomepage(api, segments)
		RegisterSegments(api, segments)

		hooks.OnStart(func() {
			slog.Info("Starting api server...")
			slog.Info(fmt.Sprintf("API: [http://%s/]", addr))
			slog.Info(fmt.Sprintf("Docs: [http://%s/docs]", addr))

			server.ListenAndServe()
		})
		// graceful shutdown
		hooks.OnStop(func() {
			slog.Info("Stopping api server...")
			ctx, cancel := context.WithTimeout(ctx, shutdownDelay)
			defer cancel()
			server.Shutdown(ctx)
		})

	})

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

	return
}

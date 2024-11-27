package lib

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/ministryofjustice/opg-reports/info"
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

// ApiTitle generates the title for the api using build details
func ApiTitle() string {
	return fmt.Sprintf("%s Reports API", info.Organisation)
}

// ApiVersion generates the version string for the api from build details.
// Currently:
//
//	`<semver> [<timestamp>] (<git-sha>)`
func ApiVersion() string {
	return info.BuildInfo()
}

// CLI returns the api wrapped as a cli command and appends not only the api routes
// but also a command to output the api spec (`openapi`)
func CLI(ctx context.Context, api huma.API, server *http.Server, segments map[string]*ApiSegment, dbPath string) (cli humacli.CLI) {
	var shutdownDelay time.Duration = 5 * time.Second

	cli = humacli.New(func(hooks humacli.Hooks, opts *CliOptions) {
		var addr = server.Addr

		AddMiddleware(api, segments, dbPath)
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

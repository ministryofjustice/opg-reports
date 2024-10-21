package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/ministryofjustice/opg-reports/versions"
)

type HomepageResponse struct {
	Body struct {
		Message string `json:"message" example:"Simple message for confirmining a connection."`
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

func main() {

	cli := humacli.New(func(hooks humacli.Hooks, opts *Opts) {
		// output info that the server is starting
		fmt.Printf("Starting api server \n  debug: \t%v\n  host: \t%s\n  port: \t%d\n\nSee [http://%s:%d/docs] for openapi spec.\n",
			opts.Debug, opts.Host, opts.Port, opts.Host, opts.Port)

		// create the server
		mux := http.NewServeMux()
		server := http.Server{
			Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
			Handler: mux,
		}

		// create the api
		versionStr := fmt.Sprintf("%s [%s] (%s)", versions.Build, versions.Timestamp, versions.Commit)
		api := humago.New(mux, huma.DefaultConfig("Reporting API", versionStr))

		// output the schema if set
		if opts.Spec {
			bytes, _ := api.OpenAPI().YAML()
			fmt.Println(string(bytes))
		}

		// register homepage action that will return an almost empty result
		huma.Register(api, huma.Operation{
			OperationID: "get-homepage",
			Method:      http.MethodGet,
			Path:        "/",
			Description: "Simple result for easier testing of connections",
		}, func(ctx context.Context, input *struct{}) (homepage *HomepageResponse, err error) {
			homepage = &HomepageResponse{}
			homepage.Body.Message = "Successful connection"
			return
		})

		// run the server
		hooks.OnStart(func() {
			server.ListenAndServe()
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

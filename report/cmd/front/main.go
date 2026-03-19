package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"opg-reports/report/internal/config"
	"opg-reports/report/packages/httpx"
	"opg-reports/report/packages/slogx"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfg  *config.Config           // the main config
	tmpl *template.Template = nil // api has an empty template
)

// main root command
var root *cobra.Command = &cobra.Command{
	Use:   "front",
	Short: `front starts the main api handler command.`,
	RunE:  runCmd,
}

// runCmd is the main cobra command
func runCmd(cmd *cobra.Command, args []string) (err error) {
	var (
		server *http.Server
		ctx    context.Context = cmd.Context()
		log    slogx.Logger    = slogx.FromContext(ctx)
		mux    httpx.Mux       = httpx.NewMux(ctx, cfg, tmpl)
	)

	// create the server
	server = &http.Server{
		Addr:    cfg.FrontHostname(),
		Handler: mux,
	}

	// attach endpoints
	registerEndpoints(ctx, mux, cfg)

	log.Info(ctx, `starting front server ...`,
		"api host", fmt.Sprintf(`http://%s`, cfg.ApiHostname()),
		"front host", fmt.Sprintf(`http://%s`, cfg.FrontHostname()),
	)
	defer log.Info(ctx, `shutting down front server ...`)

	// start the server
	server.ListenAndServe()
	return
}

func init() {
	// setup the config and bind to the command
	cfg = config.NewFront()
	cfg.BindFront(root)
}

func main() {
	var (
		err error
		ctx context.Context = context.Background()
		log slogx.Logger    = slogx.New(slogx.Config())
	)
	// attach the log & config to the current context
	ctx = slogx.Attach(ctx, log)
	// run root command
	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error(ctx, "error running command", "err", err.Error())
		os.Exit(1)
	}
}

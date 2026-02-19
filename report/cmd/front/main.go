package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/env"
	"opg-reports/report/package/logger"
	"os"

	"github.com/spf13/cobra"
)

type cli struct {
	FrontHost    string `json:"front"`         // --front-host ; this servers address
	ApiHost      string `json:"api"`           // --api-host ; this is apis address
	Version      string `json:"version"`       // --version ; the semver tag, used as part of signature
	SHA          string `json:"sha"`           // --sha ; the git commit sha used as part of signature
	RootDir      string `json:"root_dir"`      // --root-dir
	GovUKVersion string `json:"govuk_version"` // --govuk_version

	GovUKDir       string `json:"govuk_dir"`        // --govuk-dir
	LocalAssetsDir string `json:"local_assets_dir"` // --local-assets-dir
	TemplateDir    string `json:"template_dir"`     // --template-dir
	//

}

// default values for the args
var flags = &cli{
	FrontHost: ":8080",
	ApiHost:   ":8081",
	Version:   "v0.0.0",
	SHA:       "abcde",
}

// main root command
var root *cobra.Command = &cobra.Command{
	Use:   "front",
	Short: `start the front end server`,
	RunE:  runFront,
}

func registerEndpoints(ctx context.Context, mux *http.ServeMux, in *cli) {

}

// runAPI the main run command
func runFront(cmd *cobra.Command, args []string) (err error) {
	var (
		mux    *http.ServeMux
		server *http.Server
		ctx    context.Context = cmd.Context()
		log    *slog.Logger    = cntxt.GetLogger(ctx)
	)
	// overwrite arg flags from env values
	if err = env.OverwriteStruct(&flags); err != nil {
		return
	}

	// setup mux & server
	mux = http.NewServeMux()
	server = &http.Server{Addr: flags.FrontHost, Handler: mux}
	// attach endpoints
	registerEndpoints(ctx, mux, flags)
	// server info
	log.Info("Starting server ...")
	log.Info(fmt.Sprintf("VERSION: [%s] [%s]", flags.Version, flags.SHA))
	log.Info(fmt.Sprintf("API: [http://%s/]", flags.ApiHost))
	log.Info(fmt.Sprintf("FRONT: [http://%s/]", flags.FrontHost))
	// boot server
	server.ListenAndServe()
	return
}

func init() {
	root.PersistentFlags().StringVar(&flags.FrontHost, "front-host", flags.FrontHost, "Address to run this front from")
	root.PersistentFlags().StringVar(&flags.ApiHost, "api-host", flags.ApiHost, "Address of the api")
	root.PersistentFlags().StringVar(&flags.Version, "version", flags.Version, "The semver")
	root.PersistentFlags().StringVar(&flags.SHA, "sha", flags.SHA, "The git commit sha")
}

func main() {
	var err error
	var log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))
	var ctx = cntxt.AddLogger(context.Background(), log)

	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}

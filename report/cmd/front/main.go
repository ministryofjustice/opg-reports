package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/cost/costfront/home/costsbyteam"
	"opg-reports/report/internal/cost/costfront/home/costsdetailed"
	"opg-reports/report/internal/cost/costfront/home/costsdiff"
	"opg-reports/report/internal/front/homepage"
	"opg-reports/report/internal/front/statics"
	"opg-reports/report/internal/front/teampage"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/env"
	"opg-reports/report/package/logger"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type cli struct {
	FrontHost    string `json:"front"`         // --front-host ; this servers address
	ApiHost      string `json:"api"`           // --api-host ; this is apis address
	Version      string `json:"version"`       // --version ; the semver tag, used as part of signature
	SHA          string `json:"sha"`           // --sha ; the git commit sha used as part of signature
	RootDir      string `json:"root_dir"`      // --root-dir
	GovUKVersion string `json:"govuk_version"` // --govuk_version
	// fixed, based on root
	GovUKDir       string `json:"govuk_dir"`
	LocalAssetsDir string `json:"local_assets_dir"`
	TemplateDir    string `json:"template_dir"`
}

// default values for the args
var flags = &cli{
	FrontHost:      ":8080",
	ApiHost:        ":8081",
	Version:        "v0.0.0",
	SHA:            "abcde",
	RootDir:        "./",
	GovUKVersion:   "5.14.0",
	LocalAssetsDir: "web",
	TemplateDir:    "templates",
	GovUKDir:       "govuk",
}

// main root command
var root *cobra.Command = &cobra.Command{
	Use:   "front",
	Short: `start the front end server`,
	RunE:  runFront,
}

// registerEndpoints all of the front end endpoints
func registerEndpoints(ctx context.Context, mux *http.ServeMux, in *cli) {

	statics.Register(ctx, mux, &statics.Args{
		RootDir:        in.RootDir,
		GovUKDir:       in.GovUKDir,
		LocalAssetsDir: in.LocalAssetsDir,
	})
	// home pages
	// - main home
	homepage.Register(ctx, mux, &homepage.Args{
		ApiHost:      in.ApiHost,
		GovUKVersion: in.GovUKVersion,
		RootDir:      in.RootDir,
		TemplateDir:  in.TemplateDir,
	})
	// - home, simple costs grouped by team
	costsbyteam.Register(ctx, mux, &costsbyteam.Args{
		ApiHost:      in.ApiHost,
		GovUKVersion: in.GovUKVersion,
		RootDir:      in.RootDir,
		TemplateDir:  in.TemplateDir,
	})
	// - home, show detailed breakdown of costs
	costsdetailed.Register(ctx, mux, &costsdetailed.Args{
		ApiHost:      in.ApiHost,
		GovUKVersion: in.GovUKVersion,
		RootDir:      in.RootDir,
		TemplateDir:  in.TemplateDir,
	})
	// - home, show cost differences between all accounts
	costsdiff.Register(ctx, mux, &costsdiff.Args{
		ApiHost:      in.ApiHost,
		GovUKVersion: in.GovUKVersion,
		RootDir:      in.RootDir,
		TemplateDir:  in.TemplateDir,
	})

	// team pages
	// main team page
	teampage.Register(ctx, mux, &teampage.Args{
		ApiHost:      in.ApiHost,
		GovUKVersion: in.GovUKVersion,
		RootDir:      in.RootDir,
		TemplateDir:  in.TemplateDir,
	})

}

func appendRoot(in *cli) {
	in.GovUKDir = filepath.Join(in.RootDir, in.GovUKDir)
	in.LocalAssetsDir = filepath.Join(in.RootDir, in.LocalAssetsDir)
	in.TemplateDir = filepath.Join(in.RootDir, in.TemplateDir)
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
	// fix directories
	appendRoot(flags)
	// setup mux & server
	mux = http.NewServeMux()
	server = &http.Server{Addr: flags.FrontHost, Handler: mux}
	// attach endpoints
	registerEndpoints(ctx, mux, flags)
	// server info
	log.Info(fmt.Sprintf("Starting server [%s] [%s]...", flags.Version, flags.SHA))
	log.Info("Directories:")
	log.Info(fmt.Sprintf(" Root: %s", flags.RootDir))
	log.Info(fmt.Sprintf(" GovUK: %s", flags.GovUKDir))
	log.Info(fmt.Sprintf(" Local assets: %s", flags.LocalAssetsDir))
	log.Info(fmt.Sprintf(" Templates: %s", flags.TemplateDir))
	log.Info("Hosts:")
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
	root.PersistentFlags().StringVar(&flags.GovUKVersion, "govuk-version", flags.GovUKVersion, "GovUK version tag")
	root.PersistentFlags().StringVar(&flags.RootDir, "root-dir", flags.RootDir, "Root directory")
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

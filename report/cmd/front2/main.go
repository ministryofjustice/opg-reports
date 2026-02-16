package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/front/homepage"
	"opg-reports/report/internal/front/statics"
	"opg-reports/report/internal/front/teampage"
	"opg-reports/report/internal/utils/env"
	"opg-reports/report/internal/utils/logger"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	name  string = "front"
	short string = `front end server.`
)

// config items
var (
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

// used for calling api endpoints async
type blockF func(i ...any)

type cli struct {
	Name string `json:"name"`
	// data base
	DB     string `json:"db"`     // represents --db
	Driver string `json:"driver"` // represents --driver
	//
	ApiAddress string `json:"api"`     // --api
	Address    string `json:"address"` // --address
	// directories
	RootDir        string `json:"root_dir"`         // --root-dir
	GovUKDir       string `json:"govuk_dir"`        // --govuk-dir
	LocalAssetsDir string `json:"local_assets_dir"` // --local-assets-dir
	TemplateDir    string `json:"template_dir"`     // --template-dir
	//
	GovUKVersion string `json:"govuk_version"` // --govuk_version
	Signature    string `json:"signature"`     // --signature
}

// setup defaults
var flags *cli = &cli{
	Name:           "OPG Reports",
	Address:        ":8080",
	ApiAddress:     "localhost:8081",
	Driver:         "sqlite3",
	DB:             "./database/api.db",
	RootDir:        "./",
	GovUKDir:       "govuk",
	LocalAssetsDir: "local-assets",
	TemplateDir:    "templates",
	GovUKVersion:   "5.11.0",
	Signature:      "0.0.1 (abcde)",
}

var cmd = &cobra.Command{
	Use:   name,
	Short: short,
	RunE:  runFrontServer,
}

func runFrontServer(cmd *cobra.Command, args []string) (err error) {
	var (
		mux                 = http.NewServeMux()
		server              = &http.Server{Addr: flags.Address, Handler: mux}
		lg     *slog.Logger = log.With("func", "front.runFrontServer")
	)

	// update the struct from the env
	err = env.OverwriteStruct(&flags)
	if err != nil {
		return
	}
	// add root dir to the other dirs
	appendRoot()
	// static assets
	statics.Register(ctx, log, mux, &statics.Conf{
		RootDir:        flags.RootDir,
		GovUKDir:       flags.GovUKDir,
		LocalAssetsDir: flags.LocalAssetsDir,
		TemplateDir:    flags.TemplateDir,
	})

	// homepage
	homepage.Register(ctx, log, mux, &homepage.Conf{
		Name:         flags.Name,
		ApiHost:      flags.ApiAddress,
		TemplateDir:  flags.TemplateDir,
		GovUKVersion: flags.GovUKVersion,
		Signature:    flags.Signature,
	})

	teampage.Register(ctx, log, mux, &teampage.Conf{
		Name:         flags.Name,
		ApiHost:      flags.ApiAddress,
		TemplateDir:  flags.TemplateDir,
		GovUKVersion: flags.GovUKVersion,
		Signature:    flags.Signature,
	})
	// start the server
	bootInfo(lg)
	server.ListenAndServe()

	return
}

func appendRoot() {
	flags.GovUKDir = filepath.Join(flags.RootDir, flags.GovUKDir)
	flags.LocalAssetsDir = filepath.Join(flags.RootDir, flags.LocalAssetsDir)
	flags.TemplateDir = filepath.Join(flags.RootDir, flags.TemplateDir)
}
func bootInfo(lg *slog.Logger) {
	lg.Info("Starting front server...")
	lg.Info(fmt.Sprintf("Root dir: %s", flags.RootDir))
	lg.Info(fmt.Sprintf("GovUK dir: %s", flags.GovUKDir))
	lg.Info(fmt.Sprintf("Local asset dir: %s", flags.LocalAssetsDir))
	lg.Info(fmt.Sprintf("Template dir: %s", flags.TemplateDir))

	lg.Info(fmt.Sprintf("DB: %s", flags.DB))
	lg.Info(fmt.Sprintf("API: http://%s/", flags.ApiAddress))
	lg.Info(fmt.Sprintf("API Docs: http://%s/docs", flags.ApiAddress))
	lg.Info(fmt.Sprintf("Front: http://%s/", flags.Address))
}

// setup default values for config and logging
func init() {
	ctx = context.Background()
	log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))

	cmd.Flags().StringVar(&flags.DB, "db", flags.DB, "Database path")
	cmd.Flags().StringVar(&flags.Driver, "driver", flags.Driver, "Database driver")

	cmd.Flags().StringVar(&flags.Address, "address", flags.Address, "This server address")
	cmd.Flags().StringVar(&flags.ApiAddress, "api", flags.ApiAddress, "API address")

	cmd.Flags().StringVar(&flags.GovUKVersion, "govuk-version", flags.GovUKVersion, "GovUK version tag")
	cmd.Flags().StringVar(&flags.Signature, "signature", flags.Signature, "Semver signature")

	cmd.Flags().StringVar(&flags.RootDir, "root-dir", flags.RootDir, "Root directory")
	cmd.Flags().StringVar(&flags.GovUKDir, "govuk-dir", flags.GovUKDir, "GovUK directory")
	cmd.Flags().StringVar(&flags.LocalAssetsDir, "local-assets-dir", flags.LocalAssetsDir, "Local Assets directory")
	cmd.Flags().StringVar(&flags.TemplateDir, "template-dir", flags.TemplateDir, "Template directory")

}

func main() {
	var err error

	err = cmd.ExecuteContext(ctx)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}

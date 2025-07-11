package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/restr"
	"opg-reports/report/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	conf      *config.Config
	viperConf *viper.Viper
	ctx       context.Context
	log       *slog.Logger
	Info      *FrontInfo
)

type FrontInfo struct {
	AssetRoot     string            // AssetRoot is the filesystem folder that all assets sit within
	GovUKAssetDir string            // GovUKAssetDir sits under the AssetRoot nad contains downloaded items from gvuk-frontend
	LocalAssetDir string            // LocalAssetDir stores overwrites and assets custom to this project
	TemplateDir   string            // TemplateDir contains all of our templates (using .html files)
	RestClient    *restr.Repository // RestClient is used to make calls to the api via a service
}

var ()

// root command
var rootCmd = &cobra.Command{
	Use:   "front",
	Short: "front end runner",
	Long: `
front turns on the front end server.

env values that can be adjusted:

	SERVERS_FRONT_ADDR
		The address of this front end server (eg: localhost:8080)
	SERVERS_API_ADDR
		The address of the API server to connect to (eg: localhost:8081)
	GOVUK_FRONT_DIRECTORY
		The directory path to place and read govuk assets from
.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			addr   = conf.Servers.Front.Addr
			mux    = http.NewServeMux()
			server = &http.Server{Addr: addr, Handler: mux}
		)
		// Register handlers
		RegisterStaticHandlers(ctx, log, conf, Info, mux)
		RegisterHomepageHandlers(ctx, log, conf, Info, mux)

		log.Info("Starting front server...")
		log.Info(fmt.Sprintf("FRONT: [http://%s/]", addr))

		server.ListenAndServe()

		return
	},
}

// init called to setup conf / logger and to check and download
// govuk front end
func init() {
	var govdir string
	conf, viperConf = config.New()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)
	ctx = context.Background()
	govdir = filepath.Clean(conf.GovUK.Front.Directory)

	Info = &FrontInfo{
		GovUKAssetDir: govdir,
		AssetRoot:     filepath.Dir(govdir),
		LocalAssetDir: filepath.Join(filepath.Dir(govdir), "local-assets"),
		TemplateDir:   filepath.Join(filepath.Dir(govdir), "templates"),
		RestClient:    restr.Default(ctx, log, conf),
	}
	log.Info(fmt.Sprintf("ROOT ASSET DIR: [%s]", Info.AssetRoot))
	log.Info(fmt.Sprintf("GOVUK ASSET DIR: [%s]", Info.GovUKAssetDir))
	log.Info(fmt.Sprintf("LOCAL ASSET DIR: [%s]", Info.LocalAssetDir))
	log.Info(fmt.Sprintf("TEMPLATE DIR: [%s]", Info.TemplateDir))

	if !utils.DirExists(conf.GovUK.Front.Directory) {
		DownloadGovUKFrontEnd(ctx, log, conf, Info)
	}
}

func main() {
	initComponents(ctx, log, conf)
	rootCmd.Execute()
}

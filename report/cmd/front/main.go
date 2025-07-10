package main

import (
	"context"
	"log/slog"
	"net/http"
	"path/filepath"

	"opg-reports/report/config"
	"opg-reports/report/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	conf      *config.Config
	viperConf *viper.Viper
	ctx       context.Context
	log       *slog.Logger
)

var (
	assetRoot      string
	govUKAssetDir  string
	localAssetsDir string
	templateDir    string
)

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
	Run: func(cmd *cobra.Command, args []string) {
		var (
			addr   = conf.Servers.Front.Addr
			mux    = http.NewServeMux()
			server = &http.Server{Addr: addr, Handler: mux}
		)
		StartServer(ctx, log, conf, mux, server,
			RegisterStaticHandlers,
			RegisterHomepageHandlers,
		)
	},
}

// init called to setup conf / logger and to check and download
// govuk front end
func init() {
	conf, viperConf = config.New()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)
	ctx = context.Background()

	govUKAssetDir = filepath.Clean(conf.GovUK.Front.Directory)
	assetRoot = filepath.Dir(govUKAssetDir)
	localAssetsDir = filepath.Join(assetRoot, "local-assets")
	templateDir = filepath.Join(assetRoot, "templates")

	if !utils.DirExists(conf.GovUK.Front.Directory) {
		DownloadGovUKFrontEnd(ctx, log, conf)
	}
}

func main() {
	rootCmd.Execute()
}

package main

import (
	"context"
	"fmt"
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
	Info      *FrontInfo
)

type FrontInfo struct {
	AssetRoot      string
	GovUKAssetDir  string
	LocalAssetsDir string
	TemplateDir    string
	Teams          []string
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
		// get all teams when starting and attach the names
		teams, err := GetAPITeams(ctx, log, conf)
		if err != nil {
			return
		}
		for _, tm := range teams {
			if tm.Name != "Legacy" && tm.Name != "ORG" {
				Info.Teams = append(Info.Teams, tm.Name)
			}
		}

		StartServer(ctx, log, conf, Info, mux, server,
			RegisterStaticHandlers,
			RegisterHomepageHandlers,
		)
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
		Teams:          []string{},
		GovUKAssetDir:  govdir,
		AssetRoot:      filepath.Dir(govdir),
		LocalAssetsDir: filepath.Join(filepath.Dir(govdir), "local-assets"),
		TemplateDir:    filepath.Join(filepath.Dir(govdir), "templates"),
	}
	log.Info(fmt.Sprintf("ROOT ASSET DIR: [%s]", Info.AssetRoot))
	log.Info(fmt.Sprintf("GOVUK ASSET DIR: [%s]", Info.GovUKAssetDir))
	log.Info(fmt.Sprintf("LOCAL ASSET DIR: [%s]", Info.LocalAssetsDir))
	log.Info(fmt.Sprintf("TEMPLATE DIR: [%s]", Info.TemplateDir))

	if !utils.DirExists(conf.GovUK.Front.Directory) {
		DownloadGovUKFrontEnd(ctx, log, conf, Info)
	}
}

func main() {
	rootCmd.Execute()
}

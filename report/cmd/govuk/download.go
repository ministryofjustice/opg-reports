package main

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/config"
	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/service/front"
	"path/filepath"

	"github.com/spf13/cobra"
)

// frontCmd
var frontCmd = &cobra.Command{
	Use:   "front",
	Short: "front downloads the gov uk front end assets",
	Long: `
upload uploads a local database to the s3 bucket

env variables used that can be adjusted:

	SERVERS_FRONT_DIRECTORY
		The path to download into
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		err = DownloadGovUKFrontEnd(ctx, log, conf)
		return
	},
}

func DownloadGovUKFrontEnd(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
) (err error) {

	var (
		assetDir   = filepath.Clean(conf.Servers.Front.Directory)
		client     = githubr.DefaultClient(conf).Repositories
		store      = githubr.Default(ctx, log, conf)
		downloader = front.Default(ctx, log, conf)
	)
	defer log.Info("downloaded GOV UK front end ... ", "dir", assetDir, "err", err)

	files, _, err := downloader.DownloadGovUKFrontEnd(client, store, assetDir)
	if err != nil {
		log.Error("failed to download gov uk front end on init")
		return
	}
	if len(files) <= 0 {
		err = fmt.Errorf("gov uk download did not contain any files")
		log.Error("failed to get gov uk files")
		return
	}

	return
}

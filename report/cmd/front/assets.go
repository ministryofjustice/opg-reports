package main

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/config"
	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/service/front"
)

func DownloadGovUKFrontEnd(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	info *FrontInfo,
) (err error) {

	var (
		assetDir   = info.AssetRoot //filepath.Join(conf.GovUK.Front.Directory)
		client     = githubr.DefaultClient(conf).Repositories
		store      = githubr.Default(ctx, log, conf)
		downloader = front.Default[*struct{}](ctx, log, conf)
	)

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

package main

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// seedDB is called if the database doesnt exist on init, so creates a dummy one
func seedDB(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	var sqlStore sqlr.Writer = sqlr.Default(ctx, log, conf)
	var seedService *seed.Service = seed.Default(ctx, log, conf)
	_, err = seedService.All(sqlStore)
	return
}

// downloadLatestDB uses the config settings to create a s3 client and service to then download
// the latest database from the s3 bucket to local files and then updates the modified times.
func downloadLatestDB(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {

	var (
		stats        fs.FileInfo
		modifiedTime time.Time
		local        string
		dir, _                        = os.MkdirTemp("./", "__download-s3-*")
		client       *s3.Client       = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
		store        *awsr.Repository = awsr.Default(ctx, log, conf)
		now          time.Time        = time.Now().UTC()
		age          time.Duration    = 0 * time.Second
		path         string           = conf.Database.Path
		// maxAge time.Duration = 10 * time.Minute
	)

	defer func() {
		os.RemoveAll(dir)
		// always modifiy the timestamp on the database to reduce retries
		if utils.FileExists(conf.Database.Path) {
			os.Chtimes(conf.Database.Path, now, now)
		}
	}()

	stats, err = os.Stat(path)
	if err != nil {
		return
	}

	modifiedTime = stats.ModTime()
	age = now.Sub(modifiedTime)
	log.Info("database age ... ", "age", age)

	// return if age is less than max
	if age <= maxDatabaseAge {
		return
	}
	// try to fetch new databae from bucket
	log.Info("trying to download updated database ... ")
	local, err = store.DownloadItemFromBucket(client, conf.Aws.Buckets.DB.Name, conf.Aws.Buckets.DB.Path(), dir)
	if err != nil {
		log.Error("failed to download database")
		return
	}

	log.Info("downloaded new database to local file - moving ...")
	err = os.Rename(local, conf.Database.Path)
	if err != nil {
		return
	}
	return
}

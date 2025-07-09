package main

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/awsr"
	"opg-reports/report/internal/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	retryCounter     int = 0
	retryMaxAttempts int = 5
)

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
	defer os.RemoveAll(dir)

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

	if retryCounter >= retryMaxAttempts {
		log.Info("exceded rery limit for downloadind database ... skipping")
		return
	}

	// try to fetch new databae from bucket
	log.Info("trying to download updated database ... ")
	local, err = store.DownloadItemFromBucket(client, conf.Existing.DB.Bucket, conf.Existing.DB.Path(), dir)
	if err != nil {
		retryCounter++
		log.Error("failed to download database")
		return
	}

	log.Info("downloaded new database to local file - moving ...")
	err = os.Rename(local, conf.Database.Path)
	if err != nil {
		return
	}
	// modifiy the timestamp on the database to reduce retries
	if utils.FileExists(conf.Database.Path) {
		os.Chtimes(conf.Database.Path, now, now)
	}
	retryCounter = 0
	return
}

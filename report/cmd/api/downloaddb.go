package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/awsr"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/djherbis/times"
)

var (
	retryCounter     int = 0
	retryMaxAttempts int = 3
)

const maxDatabaseAge time.Duration = (24 * time.Hour) * 3 // max age of the database before refetching - 3 days

// downloadLatestDB uses the config settings to create a s3 client and service to then download
// the latest database from the s3 bucket to local files and then updates the modified times.
func downloadLatestDB(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {

	var (
		stats     times.Timespec
		createdAt time.Time
		local     string
		dir, _                     = os.MkdirTemp("./", "__download-s3-*")
		client    *s3.Client       = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
		store     *awsr.Repository = awsr.Default(ctx, log, conf)
		now       time.Time        = time.Now().UTC()
		age       time.Duration    = 0 * time.Second
		path      string           = conf.Database.Path
	)
	defer os.RemoveAll(dir)

	stats, err = times.Stat(path)
	if err != nil {
		return
	}

	// if os doesnt support btime, return
	if !stats.HasBirthTime() {
		log.Warn("OS does not support file creation times")
		return
	}

	createdAt = stats.BirthTime()
	age = now.Sub(createdAt)
	log.With("age", age, "created", createdAt).Info("database created at ... ")

	// if its younger than max, skip
	if age <= maxDatabaseAge {
		log.Info("database is within max age ... skipping ")
		return
	}
	// if we've tried to fetch it too many times, skip trying again
	if retryCounter >= retryMaxAttempts {
		log.Info("exceded rery limit for downloadind database ... skipping")
		return
	}

	// otherwise, download a fresh version
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

	retryCounter = 0
	return
}

package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domain/downloads/download"
	"opg-reports/report/internal/utils/awsclients"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

const (
	dlCmdName   string = "download" // root command name
	dlShortDesc string = `download is used to download the sqlite db from the configured s3 bucket.`
	dlLongDesc  string = `
download is used to download the sqlite db from the configured s3 bucket to local file system.
`
)

var (
	dlBucket    string = "opg-reports-development"
	dlKey       string = "database/api.db"
	dlDirectory string = "./dl"
	dlRegion    string = "eu-west-1"
)

var (
	dlCmd *cobra.Command = &cobra.Command{
		Use:   dlCmdName,
		Short: dlShortDesc,
		Long:  dlLongDesc,
		RunE:  dlRunE,
	}
)

// wrapper to use with cobra
func dlRunE(cmd *cobra.Command, args []string) (err error) {

	var client *s3.Client
	client, err = awsclients.New[*s3.Client](ctx, log, dlRegion)
	return downloadItem(ctx, log, client, &download.Options{
		Bucket:    dlBucket,
		Key:       dlKey,
		Directory: dlDirectory,
	})
}

func downloadItem(ctx context.Context, log *slog.Logger, client download.AwsClient, opts *download.Options) (err error) {
	var path string
	var lg *slog.Logger = log.With("func", "db.downloadItem")

	lg.Info("starting db download command ...")
	lg.With("opts", opts).Debug("options ...")

	path, err = download.DownloadItemFromBucket(ctx, log, client, opts)
	if err != nil {
		return
	}
	lg.With("path", path).Info("complete.")
	return
}

func init() {
	dlCmd.Flags().StringVar(&dlBucket, "bucket", dlBucket, "Bucket name to fetch from.")
	dlCmd.Flags().StringVar(&dlKey, "key", dlKey, "Item key to download from the bucket")
	dlCmd.Flags().StringVar(&dlDirectory, "directory", dlDirectory, "Top level directory to download into.")
	dlCmd.Flags().StringVar(&dlRegion, "region", dlRegion, "AWS region.")
}

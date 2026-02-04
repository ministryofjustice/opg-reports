package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domain/uploads/upload"
	"opg-reports/report/internal/utils/awsclients"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

const (
	upCmdName   string = "upload" // root command name
	upShortDesc string = `upload is used to upload the sqlite db to the s3 bucket.`
	upLongDesc  string = `
upload is used to upload the sqlite db local to the configured s3 bucket.
`
)

var (
	upBucket string = "opg-reports-development"
	upKey    string = "database/api.db"
	upFile   string = "database/api.db"
	upRegion string = "eu-west-1"
)

var (
	upCmd *cobra.Command = &cobra.Command{
		Use:   upCmdName,
		Short: upShortDesc,
		Long:  upLongDesc,
		RunE:  upRunE,
	}
)

// wrapper to use with cobra
func upRunE(cmd *cobra.Command, args []string) (err error) {

	var client *s3.Client
	client, err = awsclients.New[*s3.Client](ctx, log, dlRegion)
	return uploadItem(ctx, log, client, &upload.Options{
		Bucket:   upBucket,
		Key:      upKey,
		Filepath: "",
	})
}

func uploadItem(ctx context.Context, log *slog.Logger, client upload.AwsClient, opts *upload.Options) (err error) {
	var lg *slog.Logger = log.With("func", "db.uploadItem")

	lg.Info("starting db upload command ...")
	lg.With("opts", opts).Debug("options ...")

	_, err = upload.UploadItemToBucket(ctx, log, client, opts)
	if err != nil {
		return
	}

	lg.With("filepath", opts.Filepath).Info("complete.")
	return
}

func init() {
	upCmd.Flags().StringVar(&upBucket, "bucket", upBucket, "Bucket name")
	upCmd.Flags().StringVar(&upKey, "key", upKey, "Item key to upload file as.")
	upCmd.Flags().StringVar(&upFile, "file", upFile, "File to upload.")
	upCmd.Flags().StringVar(&upRegion, "region", upRegion, "AWS region.")
}

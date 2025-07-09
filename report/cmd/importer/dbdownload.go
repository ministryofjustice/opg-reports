package main

import (
	"os"

	"opg-reports/report/internal/repository/awsr"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

// dbDownloadCmd downloads the database from the s3 bucket to a temp file
// and then overwrites (using os.Rename) the configured database file.
var dbDownloadCmd = &cobra.Command{
	Use:   "dbdownload",
	Short: "dbdownload downloads the database from an s3 bucket to local file system",
	Long: `
dbdownload downloads the database from an s3 bucket to local file system

env variables used that can be adjusted:

	AWS_BUCKETS_DB_NAME
		The name of the bucket that stores the sqlite database
	AWS_BUCKETS_DB_KEY
		The object key in the bucket (include folder path) where the sqlite db is stored
	DATABASE_PATH
		The file path to the sqlite database on the local filesystem to copy the s3 version into
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			s3Client = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
			awsStore = awsr.Default(ctx, log, conf)
		)
		err = dbDownloadCmdRunner(s3Client, awsStore)
		return
	},
}

func dbDownloadCmdRunner(client awsr.ClientS3Getter, store awsr.RepositoryS3BucketItemDownloader) (err error) {
	var (
		dir, _ = os.MkdirTemp("./", "__download-s3-*")
		local  string
	)
	defer os.RemoveAll(dir)
	local, err = store.DownloadItemFromBucket(client, conf.Aws.Buckets.DB.Name, conf.Aws.Buckets.DB.Path(), dir)
	if err != nil {
		return
	}
	err = os.Rename(local, conf.Database.Path)
	return
}

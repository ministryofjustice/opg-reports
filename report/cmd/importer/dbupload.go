package main

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
)

// dbUploadCmd uploads the local database to the configured s3 bucket
var dbUploadCmd = &cobra.Command{
	Use:   "dbupload",
	Short: "dbupload uploads a local database to the s3 bucket",
	Long: `
dbupload uploads a local database to the s3 bucket

env variables used that can be adjusted:

	AWS_BUCKETS_DB_NAME
		The name of the bucket to upload the database to
	AWS_BUCKETS_DB_KEY
		The object key for the bucket (including folder path) where the sqlite db will be uploaded
	DATABASE_PATH
		The file path to the sqlite database on the local filesystem to upload to s3
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			s3Client = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
			awsStore = awsr.Default(ctx, log, conf)
		)
		err = dbUploadCmdRunner(s3Client, awsStore)
		return
	},
}

func dbUploadCmdRunner(
	client awsr.ClientS3Putter,
	store awsr.RepositoryS3BucketItemUploader,
) (err error) {
	var (
		dir, _       = os.MkdirTemp("./", "__upload-s3-*")
		copyFrom     = conf.Database.Path
		copyTo       = filepath.Join(dir, filepath.Base(conf.Database.Path))
		targetBucket = conf.Aws.Buckets.DB.Name
		targetKey    = conf.Aws.Buckets.DB.Path()
		src          *os.File
	)
	// open the existing db file & copy to the new location
	src, err = os.Open(copyFrom)
	if err != nil {
		return
	}
	defer func() {
		src.Close()
		os.RemoveAll(dir)
	}()
	// copy...
	err = utils.FileCopy(src, copyTo)
	if err != nil {
		return
	}
	// targetKey = "database/api2.db"
	// now upload the copy
	_, err = store.UploadItemToBucket(client, targetBucket, targetKey, copyTo)

	return
}

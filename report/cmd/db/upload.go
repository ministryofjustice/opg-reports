package main

import (
	"os"
	"path/filepath"

	"opg-reports/report/internal/repository/awsr"
	"opg-reports/report/internal/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

// uploadCmd uploads the local database to the configured s3 bucket
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload uploads a local database to the s3 bucket",
	Long: `
upload uploads a local database to the s3 bucket

env variables used that can be adjusted:

	DATABASE_BUCKET_NAME
		The name of the bucket to upload the database to
	DATABASE_PATH
		The file path to the sqlite database on the local filesystem to upload to s3
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			s3Client = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
			awsStore = awsr.Default(ctx, log, conf)
		)
		err = uploadCmdRunner(s3Client, awsStore)
		return
	},
}

func uploadCmdRunner(
	client awsr.ClientS3Putter,
	store awsr.RepositoryS3BucketItemUploader,
) (err error) {
	var (
		dir, _       = os.MkdirTemp("./", "__upload-s3-*")
		copyFrom     = conf.Database.Path
		copyTo       = filepath.Join(dir, filepath.Base(conf.Database.Path))
		targetBucket = conf.Database.Bucket.Name
		targetKey    = conf.Database.Bucket.Path()
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

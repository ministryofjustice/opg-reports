package awsr

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ListBucket returns all the files stored within the bucket and under the prefix
// path passed.
//
// Used to find files in remote s3 buckets by various services.
//
// client - ListObjectsV2APIClient
func (self *Repository) ListBucket(client s3.ListObjectsV2APIClient, bucket string, prefix string) (files []string, err error) {
	var (
		paginator *s3.ListObjectsV2Paginator
		log       *slog.Logger = self.log.With("operation", "ListBucket", "bucket", bucket, "prefix", prefix)
		pg        int          = 0
	)

	files = []string{}

	log.Debug("getting paginated list of files in bucket ...")
	paginator = s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{Bucket: &bucket, Prefix: &prefix})
	// loop over all the pages in the pagination set fetch them all
	for paginator.HasMorePages() {
		var page *s3.ListObjectsV2Output
		pg++
		log.With("page", pg).Debug("getting page in list ...")
		page, err = paginator.NextPage(context.Background())
		if err != nil {
			log.With("page", pg).Error("failed to fetch list")
			return
		}
		// append the file to the list
		for _, obj := range page.Contents {
			files = append(files, *obj.Key)
		}
	}

	return
}

// DownloadBucket finds and fetches all files under the <prefix> path and saves them to local file underneath <directory>
// and maintains the path used by the bucket. Returns a list of all the files with their local file paths.
func (self *Repository) DownloadBucket(client ClientS3ListAndGetter, bucket string, prefix string, directory string) (downloaded []string, err error) {
	var (
		log   *slog.Logger = self.log.With("operation", "DownloadBucket", "bucket", bucket, "prefix", prefix)
		files []string     = []string{}
	)
	log.Debug("downloading bucket to local directory ...")

	os.MkdirAll(directory, os.ModePerm)
	downloaded = []string{}
	// get all files from the bucket
	files, err = self.ListBucket(client, bucket, prefix)
	if err != nil {
		return
	}

	for _, file := range files {

		saved, e := self.DownloadItemFromBucket(client, bucket, file, directory)
		if e == nil {
			downloaded = append(downloaded, saved)
		}

	}

	return
}

// DownloadItemFromBucket fetches a single item from the s3 bucket and saves it to a path underneath <directory>
// and maintaines the bucket path for hte local file.
func (self *Repository) DownloadItemFromBucket(client ClientS3Getter, bucket string, key string, directory string) (file string, err error) {
	var (
		result    *s3.GetObjectOutput
		body      []byte
		localFile string       = filepath.Join(directory, key)
		parentDir string       = filepath.Dir(localFile)
		log       *slog.Logger = self.log.With("operation", "DownloadItemFromBucket", "bucket", bucket, "key", key)
		opts                   = &s3.GetObjectInput{Bucket: &bucket, Key: &key}
		ctx                    = self.ctx
	)
	log.Debug("downloading item from s3 ...")
	os.MkdirAll(parentDir, os.ModePerm)
	// s3
	result, err = client.GetObject(ctx, opts)
	if err != nil {
		return
	}
	defer result.Body.Close()
	body, _ = io.ReadAll(result.Body)

	if err = os.WriteFile(localFile, body, os.ModePerm); err == nil {
		file = localFile
	}
	return
}

// UploadItemToBucket uploads a local file to a bucket.
//
// Key must contain any bucket prefix / path and the complete filename. the localFile path is
// only used for reading content
//
// Uses the default sse of aes256 rather than a custom kms key
func (self *Repository) UploadItemToBucket(client ClientS3Putter, bucket string, key string, localFile string) (result *s3.PutObjectOutput, err error) {
	var (
		log *slog.Logger = self.log.With("operation", "UploadItemToBucket", "bucket", bucket, "key", key, "localFile", localFile)
		ctx              = self.ctx
	)
	log.Debug("uploading item to s3 ...")
	file, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer file.Close()

	result, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               &bucket,
		Key:                  &key,
		Body:                 file,
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	})

	return
}

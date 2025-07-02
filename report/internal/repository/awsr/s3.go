package awsr

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client fetches a v2 version of the sts client loading from the
// env and setting the region
//
// Used to establish connection to the aws api for s3 bucket calls
func ClientS3(ctx context.Context, region string) (client *s3.Client, err error) {
	var awscfg aws.Config

	awscfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return
	}
	client = s3.NewFromConfig(awscfg)
	return

}

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
//
// Uses go concurrency to fetch the files to speed things up
func (self *Repository) DownloadBucket(client ClientS3ListAndGetter, bucket string, prefix string, directory string) (downloaded []string, err error) {
	var (
		log   *slog.Logger   = self.log.With("operation", "DownloadBucket", "bucket", bucket, "prefix", prefix)
		mutex *sync.Mutex    = &sync.Mutex{}
		wg    sync.WaitGroup = sync.WaitGroup{}
		files []string       = []string{}
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
		wg.Add(1)
		//
		go func() {
			saved, e := self.DownloadItemFromBucket(client, bucket, file, directory)
			if e == nil {
				mutex.Lock()
				downloaded = append(downloaded, saved)
				mutex.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return
}

// DownloadItemFromBucket fetches a single item from the s3 bucket and saves it to a path underneath <directory>
// and maintaines the bucket path for hte local file.
//
// If <client> is nil then it will try to generate a fresh client
func (self *Repository) DownloadItemFromBucket(client ClientS3Getter, bucket string, key string, directory string) (file string, err error) {
	var (
		result    *s3.GetObjectOutput
		body      []byte
		localFile string       = filepath.Join(directory, key)
		parentDir string       = filepath.Dir(localFile)
		log       *slog.Logger = self.log.With("operation", "DownloadItemFromBucket", "bucket", bucket, "key", key)
		ctx, _                 = context.WithTimeout(self.ctx, 10*time.Second)
		opts                   = &s3.GetObjectInput{Bucket: &bucket, Key: &key}
	)
	log.Debug("downloading item from s3 ...")
	os.MkdirAll(parentDir, os.ModePerm)

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

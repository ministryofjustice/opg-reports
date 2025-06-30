package s3bucket

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ministryofjustice/opg-reports/report/config"
)

// Repository
//
// interfaces:
//   - Repository
//   - S3Repository
type Repository struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

// connection is an internal helper to handle creating the client
func (self *Repository) connection() (client *s3.S3, err error) {
	var (
		sess           *session.Session
		sessionOptions *aws.Config = &aws.Config{Region: aws.String(self.conf.Aws.GetRegion())}
	)
	// create a new session
	sess, err = session.NewSession(sessionOptions)
	if err != nil {
		return
	}
	client = s3.New(sess)

	return
}

// ListBucket returns all s3 object keys found within the bucket passed that are under the
// prefix.
//
// Used to find all files ready for downloading to local host
func (self *Repository) ListBucket(bucket string, prefix string) (fileList []string, err error) {
	var (
		client      *s3.S3
		log         = self.log.With("bucket", bucket, "prefix", prefix, "operation", "ListBucket")
		listOptions = &s3.ListObjectsInput{
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		}
	)
	log.Debug("creating s3 client ...")
	client, err = self.connection()
	if err != nil {
		return
	}

	log.Debug("listing objects in bucket ...")
	err = client.ListObjectsPagesWithContext(self.ctx, listOptions, func(o *s3.ListObjectsOutput, b bool) bool {
		for _, o := range o.Contents {
			fileList = append(fileList, *o.Key)
		}
		return true
	})
	log.With("count", len(fileList)).Debug("found objects in bucket")

	return
}

// Download fetches the set of files from the s3 bucket and write them to the locaDir while maintaining any sub folder structures.
//
// Uses a batch download object with an After hook that write the content of the buffer to the file system - this allows many
// downloads at once and no need to run save afterwards.
//
// If the number of downloadedFiles does not match the number files an error is returned.
func (self *Repository) Download(bucket string, files []string, localDir string) (downloadedFiles []string, err error) {
	var (
		client        *s3.S3
		batchDownload = []s3manager.BatchDownloadObject{}
		log           = self.log.With("bucket", bucket, "operation", "Download")
	)
	log.Debug("creating s3 client ...")
	client, err = self.connection()
	if err != nil {
		return
	}
	// make the local directory
	os.MkdirAll(localDir, os.ModePerm)
	downloadedFiles = []string{}

	log.Debug("creating batch download list ...")
	for _, file := range files {
		var (
			localFile = filepath.Join(localDir, file)
			parentDir = filepath.Dir(localFile)
		)
		os.MkdirAll(parentDir, os.ModePerm)

		buff := aws.NewWriteAtBuffer([]byte{})
		// create a batch download object with s3 info and after trigger that with write the file
		// to local files
		batchDownload = append(batchDownload, s3manager.BatchDownloadObject{
			Object: &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(file),
			},
			Writer: buff,
			After: func() error {
				log.Debug("writing downloaded file", "destination", localFile)
				if e := os.WriteFile(localFile, buff.Bytes(), os.ModePerm); e == nil {
					downloadedFiles = append(downloadedFiles, localFile)
				}
				return nil
			},
		})

	}

	// setup with interator
	downloader := s3manager.NewDownloaderWithClient(client)
	err = downloader.DownloadWithIterator(self.ctx, &s3manager.DownloadObjectsIterator{
		Objects: batchDownload,
	})

	if len(downloadedFiles) != len(files) {
		err = fmt.Errorf("downloaded a different nubmer of files, expected [%d] actual [%d]", len(files), len(downloadedFiles))
	}
	return
}

// New provides a configured repository instance
func New(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository, err error) {
	rp = &Repository{}

	if log == nil {
		err = fmt.Errorf("no logger passed for s3bucket repository")
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for s3bucket repository")
		return
	}

	log = log.WithGroup("s3bucket")
	rp = &Repository{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}

	return
}

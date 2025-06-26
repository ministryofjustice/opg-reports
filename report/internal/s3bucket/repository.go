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

// Download fetches the set of files from the s3 bucket and write them to the locaDir while maintaining any sub folder structures
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

		batchDownload = append(batchDownload, s3manager.BatchDownloadObject{
			Object: &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(file),
			},
			Writer: &buffer{Path: localFile},
		})

	}
	// run at end to write to files to local storage
	defer func() {
		for _, obj := range batchDownload {
			buf, ok := obj.Writer.(*buffer)
			lg := log.With("objectKey", *obj.Object.Key)

			if !ok {
				lg.Error("issue with downloading item")
				continue
			}
			n, err := buf.Save()
			if err != nil {
				lg.Error("error saving buffer to local file")
				continue
			}
			if n < 1 {
				lg.Error("error with saved buffer size")
				continue
			}
			lg.Warn("downloaded to " + buf.Path)
			downloadedFiles = append(downloadedFiles, buf.Path)
		}
	}()
	// handle with interator
	err = s3manager.NewDownloaderWithClient(client).DownloadWithIterator(self.ctx, &s3manager.DownloadObjectsIterator{Objects: batchDownload})
	return
}

type buffer struct {
	Path string
	buf  *aws.WriteAtBuffer
}

func (b *buffer) WriteAt(p []byte, off int64) (n int, err error) {
	if b.buf == nil {
		b.buf = &aws.WriteAtBuffer{}
	}
	return b.buf.WriteAt(p, off)
}

func (b *buffer) Save() (n int, err error) {
	if b.buf == nil {
		return
	}
	var f *os.File
	f, err = os.Create(b.Path)
	if err != nil {
		return
	}
	defer f.Close()
	n, err = f.Write(b.buf.Bytes())
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

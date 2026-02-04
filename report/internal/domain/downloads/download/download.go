package download

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	ErrFailedToGetObject  = errors.New("failed to get object from s3.")
	ErrFailedToReadObject = errors.New("failed to read object content.")
	ErrFailedToWriteFile  = errors.New("failed to write content to local file.")
)

// AwsClient is used to allow mocking and is a proxy for *s3.Client
type AwsClient interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type Options struct {
	Bucket    string
	Key       string
	Directory string // local storage location to write file to
}

// DownloadItemFromBucket connects to s3 setup and pull the content of a single item from s3 - presumes the location is know nexactly.
func DownloadItemFromBucket[T AwsClient](ctx context.Context, log *slog.Logger, client T, options *Options) (path string, err error) {
	var (
		result    *s3.GetObjectOutput
		body      []byte
		localFile string             = filepath.Join(options.Directory, options.Key)
		localDir  string             = filepath.Dir(localFile)
		lg        *slog.Logger       = log.With("func", "domain.downloads.download.GetItemFromBucket")
		opts      *s3.GetObjectInput = &s3.GetObjectInput{
			Bucket: &options.Bucket,
			Key:    &options.Key,
		}
	)

	lg.Debug("starting ...")
	os.MkdirAll(localDir, os.ModePerm)
	// get the object from s3
	lg.With("options", options).Debug("getting object ...")
	result, err = client.GetObject(ctx, opts)
	if err != nil {
		lg.Error("error getting object", "err", err.Error())
		err = errors.Join(ErrFailedToGetObject, err)
		return
	}
	defer result.Body.Close()
	// read the content
	body, err = io.ReadAll(result.Body)
	if err != nil {
		lg.Error("error reading object", "err", err.Error())
		err = errors.Join(ErrFailedToReadObject, err)
		return
	}
	// write content to local disk
	err = os.WriteFile(localFile, body, os.ModePerm)
	if err != nil {
		lg.Error("error writing file", "err", err.Error(), "file", localFile)
		err = errors.Join(ErrFailedToReadObject, err)
		return
	}
	path = localFile
	lg.Debug("complete.")
	return
}

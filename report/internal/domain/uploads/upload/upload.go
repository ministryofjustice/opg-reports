package upload

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	ErrFailedToGetObject  = errors.New("failed to get object from s3.")
	ErrFailedToReadObject = errors.New("failed to read object content.")
	ErrFailedToWriteFile  = errors.New("failed to write content to local file.")
)

// AwsClient is used to allow mocking and is a proxy for *s3.Client
type AwsClient interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type Options struct {
	Bucket   string
	Key      string
	Filepath string //
}

// UploadItemToBucket uploads a local file to a bucket.
//
// Key must contain any bucket prefix / path and the complete filename. the localFile path is
// only used for reading content
//
// Uses the default sse of aes256 rather than a custom kms key
func UploadItemToBucket[T AwsClient](ctx context.Context, log *slog.Logger, client T, options *Options) (result *s3.PutObjectOutput, err error) {
	var (
		file *os.File
		opts *s3.PutObjectInput
		lg   *slog.Logger = log.With("func", "download.GetItemFromBucket")
	)

	lg.Debug("starting ...")

	lg.With("localfile", options.Filepath).Debug("opening file ...")
	file, err = os.Open(options.Filepath)
	if err != nil {
		return
	}
	defer file.Close()

	lg.Debug("putting object ...")
	opts = &s3.PutObjectInput{
		Bucket:               &options.Bucket,
		Key:                  &options.Key,
		Body:                 file,
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	}
	result, err = client.PutObject(ctx, opts)
	if err != nil {
		return
	}

	lg.Debug("complete.")
	return
}

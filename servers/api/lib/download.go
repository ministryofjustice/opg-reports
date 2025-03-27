package lib

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ministryofjustice/opg-reports/internal/awscfg"
	"github.com/ministryofjustice/opg-reports/internal/awssession"
)

// DownloadS3DB tries to download the bucketDB from the bucketName and write the content
// to localPath.
// It will write the content to a temp file in the same directory first, and then rename
// the file on success
func DownloadS3DB(bucketName string, bucketDB string, localPath string) (ok bool, err error) {
	var (
		sess         *session.Session
		svc          *s3.S3
		result       *s3.GetObjectOutput
		tempFile     *os.File
		body         []byte
		tempFilename string
		awsCfg       *awscfg.Config = awscfg.FromEnv()
		localDir     string         = filepath.Dir(localPath)
		ext          string         = filepath.Ext(localPath)
	)
	ok = true

	if sess, err = awssession.New(awsCfg); err != nil {
		ok = false
		slog.Error("[api] downloading from s3 - aws session failed", slog.String("err", err.Error()))
		return
	}
	svc = s3.New(sess)
	// // use a head object call to see if file exists
	// _, err = svc.HeadObject(&s3.HeadObjectInput{
	// 	Bucket: aws.String(bucketName),
	// 	Key:    aws.String(bucketDB),
	// })
	// // if the head call failed, return error
	// if err != nil {
	// 	ok = false
	// 	slog.Error("[api] s3 head call failed, object doesnt exist", slog.String("err", err.Error()))
	// 	return
	// }

	// get the object from the bucket
	result, err = svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketDB),
	})
	if err != nil {
		ok = false
		slog.Error("[api] s3 get object failed", slog.String("err", err.Error()))
		return
	}
	defer result.Body.Close()

	// make the dir
	os.MkdirAll(localDir, os.ModePerm)
	// make the file
	tempFile, err = os.CreateTemp(localDir, fmt.Sprintf("*%s", ext))
	if err != nil {
		ok = false
		slog.Error("[api] temp file creation failed", slog.String("err", err.Error()))
		return
	}
	defer tempFile.Close()
	tempFilename = tempFile.Name()

	body, err = io.ReadAll(result.Body)
	if err != nil {
		ok = false
		slog.Error("[api] reading result body failed", slog.String("err", err.Error()))
		return
	}

	_, err = tempFile.Write(body)
	if err == nil {
		slog.Info("[api] downloaded file, moving", slog.String("old", tempFilename), slog.String("new", localPath))
		os.Remove(localPath)
		os.Rename(tempFilename, localPath)
	}

	return
}

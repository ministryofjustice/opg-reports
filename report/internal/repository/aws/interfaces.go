package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type S3er interface {
	ListBucket(bucket string, prefix string) (files []string, err error)
	DownloadBucket(bucket string, prefix string, directory string) (downloaded []string, err error)
	DownloadItemFromBucket(bucket string, key string, directory string, client *s3.Client) (file string, err error)
}

type STSer interface {
	GetCallerIdentity() (caller *sts.GetCallerIdentityOutput, err error)
}

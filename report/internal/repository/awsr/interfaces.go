package awsr

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Model interface{}

type STSer interface {
	GetCallerIdentity(client ClientSTSCaller) (caller *sts.GetCallerIdentityOutput, err error)
}

type S3er interface {
	ListBucket(client s3.ListObjectsV2APIClient, bucket string, prefix string) (files []string, err error)
	DownloadBucket(client ClientS3ListAndGetter, bucket string, prefix string, directory string) (downloaded []string, err error)
	DownloadItemFromBucket(client ClientS3Getter, bucket string, key string, directory string) (file string, err error)
}

type ClientSTSCaller interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

type ClientS3Getter interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type ClientS3ListAndGetter interface {
	s3.ListObjectsV2APIClient
	ClientS3Getter
}

type ClientS3 interface {
	ClientS3ListAndGetter
}

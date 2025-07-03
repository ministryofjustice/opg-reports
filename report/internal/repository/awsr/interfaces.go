package awsr

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Model interface{}

// STSRepository is interface represeting the functionality of this repository module that relates
// to STS part of AWS SDK
type STSRepository interface {
	GetCallerIdentity(client ClientSTSCaller) (caller *sts.GetCallerIdentityOutput, err error)
}

// S3er is an interface that represents the functionality of this repository module which relates
// to S3 - so listing & downloading mostly
type S3Repository interface {
	ListBucket(client s3.ListObjectsV2APIClient, bucket string, prefix string) (files []string, err error)
	DownloadBucket(client ClientS3ListAndGetter, bucket string, prefix string, directory string) (downloaded []string, err error)
	DownloadItemFromBucket(client ClientS3Getter, bucket string, key string, directory string) (file string, err error)
}

// ClientSTSCaller represents the client (sts.Client) interface used by STSer
// to access the SDK.
//
// Used to allow for mocking
type ClientSTSCaller interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

// ClientS3Getter represents the client (s3.Client) used to download an item
// from a bucket.
//
// Allows mocking with custom structs
type ClientS3Getter interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// ClientS3ListAndGetter represents both listing and downloading capabilities
type ClientS3ListAndGetter interface {
	s3.ListObjectsV2APIClient
	ClientS3Getter
}

// ClientS3 represents an overal S3 client that could be mocked
type ClientS3 interface {
	ClientS3ListAndGetter
}

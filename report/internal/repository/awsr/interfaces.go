package awsr

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Model interface{}

// RepositorySTS is interface represeting the functionality of this repository module that relates
// to STS part of AWS SDK
type RepositorySTS interface {
	GetCallerIdentity(client ClientSTSCaller) (caller *sts.GetCallerIdentityOutput, err error)
}

// RepositoryS3 is an interface that represents the functionality of this repository module which relates
// to S3 - so listing & downloading mostly
type RepositoryS3 interface {
	RepositoryS3BucketLister
	RepositoryS3BucketDownloader
	RepositoryS3BucketItemDownloader
}

// RepositoryS3BucketLister interface requires the repository provides a way to list the content of a bucket.
type RepositoryS3BucketLister interface {
	ListBucket(client s3.ListObjectsV2APIClient, bucket string, prefix string) (files []string, err error)
}

// RepositoryS3BucketDownloader requires the repository has a method to download all files in the bucket under the prefix
// to the local file system
type RepositoryS3BucketDownloader interface {
	DownloadBucket(client ClientS3ListAndGetter, bucket string, prefix string, directory string) (downloaded []string, err error)
}

// RepositoryS3BucketItemDownloader requires the repository has a function to download a specific item from a bucket to local
// file system
type RepositoryS3BucketItemDownloader interface {
	DownloadItemFromBucket(client ClientS3Getter, bucket string, key string, directory string) (file string, err error)
}
type RepositoryS3BucketItemUploader interface {
	UploadItemToBucket(client ClientS3Putter, bucket string, key string, localFile string) (result *s3.PutObjectOutput, err error)
}

// RepositoryCostExplorer contains all the methods used to fetch cost data from the aws sdk
type RepositoryCostExplorer interface {
	RepositoryCostExplorerGetter
}

// RepositoryCostExplorerGetter provides all method to get cost and usage data from the aws sdk
type RepositoryCostExplorerGetter interface {
	GetCostData(client ClientCostExplorerGetter, options *GetCostDataOptions) (values []map[string]string, err error)
}

// ClientSTSCaller represents the client (sts.Client) interface used by RepositorySTS
// to access the SDK.
//
// Used to allow for mocking
type ClientSTSCaller interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

// ClientS3 represents an overal S3 client
type ClientS3 interface {
	ClientS3ListAndGetter
	ClientS3Putter
}

// ClientS3Getter represents the client (s3.Client) used to download an item
// from a bucket.
type ClientS3Getter interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}
type ClientS3Putter interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// ClientS3ListAndGetter represents both listing and downloading capabilities that a client
// needs to be able to find and download all items from a bucket
type ClientS3ListAndGetter interface {
	s3.ListObjectsV2APIClient
	ClientS3Getter
}

// ClientCostExplorerGetter represents the method needed by a client (costexplorer.Client) to call
// the aws sdk
type ClientCostExplorerGetter interface {
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

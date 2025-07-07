package awsr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// check interfaces are correct
var (
	_ S3Repository  = &Repository{}
	_ STSRepository = &Repository{}
)

var accountId = "001A"
var accountArn = "arn:aws:iam::001A:user/person/test.name"
var accountUser = "AIDZFKC"

// mockSTSCaller used to avoid making api calls
type mockSTSCaller struct{}

// GetCallerIdentity always returns dummy creds
func (self *mockSTSCaller) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return &sts.GetCallerIdentityOutput{
		Account: &accountId,
		Arn:     &accountArn,
		UserId:  &accountUser,
	}, nil
}

type mockS3BucketLister struct{}

func (self *mockS3BucketLister) ListBucket(client s3.ListObjectsV2APIClient, bucket string, prefix string) ([]string, error) {
	return []string{
		fmt.Sprintf("%s/%s%s", bucket, prefix, "sample-00.json"),
		fmt.Sprintf("%s/%s%s", bucket, prefix, "sample-01.json"),
		fmt.Sprintf("%s/%s%s", bucket, prefix, "sample-01.csv"),
	}, nil
}

// mockedS3BucketDownloader provides a mocked version of DownloadBucket that writes a dummy cost file to a
// known location and returns that as the file path
type mockedS3BucketDownloader struct{}

// DownloadBucket generates a file with dummy cost data in to for testing inserts
func (self *mockedS3BucketDownloader) DownloadBucket(client ClientS3ListAndGetter, bucket string, prefix string, directory string) (downloaded []string, err error) {
	var file = filepath.Join(directory, "sample-costs.json")
	var content = `[
	{
		"id": 0,
		"ts": "2024-08-15 18:52:55.055478 +0000 UTC",
		"organisation": "OPG",
		"account_id": "001A",
		"account_name": "Account 1A",
		"unit": "TEAM-A",
		"label": "A",
		"environment": "development",
		"service": "Amazon Simple Storage Service",
		"region": "eu-west-1",
		"date": "2025-05-31",
		"cost": "0.2309542206"
	},
	{
		"id": 0,
		"ts": "2024-08-15 18:52:55.055478 +0000 UTC",
		"organisation": "OPG",
		"account_id": "001A",
		"account_name": "Account 1A",
		"unit": "TEAM-A",
		"label": "A",
		"environment": "development",
		"service": "Amazon Simple Storage Service",
		"region": "eu-west-1",
		"date": "2025-04-31",
		"cost": "107.53"
	}
]`
	err = os.WriteFile(file, []byte(content), os.ModePerm)
	downloaded = append(downloaded, file)
	return
}

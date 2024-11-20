/*
convertor takes an older formatted data file and converts over to new data file
*/
package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	v1 "github.com/ministryofjustice/opg-reports/convertor/v1"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/pkg/awscfg"
	"github.com/ministryofjustice/opg-reports/pkg/awssession"
)

const (
	bucketName string = "report-data-development"
)

var (
	awsCfg       = awscfg.FromEnv()
	dataDir      = "./data"
	convertedDir = "./converted"
)

// Download grabs all the files from the s3 bucket and clones them locally
func Download(sess *session.Session, svc *s3.S3) {
	fmt.Printf("Downloading from s3 bucket [%s]\n", bucketName)
	waitgroup := sync.WaitGroup{}
	downloader := s3manager.NewDownloader(sess)
	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		err = fmt.Errorf("Unable to list items in bucket %q, %v", bucketName, err)
		return
	}
	// loop over all files and get the details
	items := []string{}
	more := true
	for more {
		for _, item := range resp.Contents {
			items = append(items, *item.Key)
		}
		more = *resp.IsTruncated
		if more {
			resp, _ = svc.ListObjectsV2(&s3.ListObjectsV2Input{
				Bucket:            aws.String(bucketName),
				ContinuationToken: resp.NextContinuationToken,
			})
		}
	}

	// no loop over the file items and download them all
	for _, obj := range items {
		waitgroup.Add(1)
		go func(item string) {
			var (
				file      *os.File
				key       *string = aws.String(item)
				bucketDir string  = filepath.Dir(item)
				dir       string  = filepath.Join(dataDir, bucketDir)
				path      string  = filepath.Join(dataDir, item)
			)

			os.MkdirAll(dir, os.ModePerm)
			file, err = os.Create(path)
			if err != nil {
				panic(err)
			}

			_, err = downloader.Download(file, &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    key,
			})
			file.Close()
			if err != nil {
				panic(err)
			}

			waitgroup.Done()
		}(obj)
	}
	waitgroup.Wait()
	fmt.Println("Downloaded.")
}

// ConvertV1s takes the known sub dirs in the bucket thats been cloned locally
// and converts the older structs to new ones
func ConvertV1s() {
	fmt.Println("Converting v1s ...")
	var (
		costs     = []*models.AwsCost{}
		uptimes   = []*models.AwsUptime{}
		standards = []*models.GitHubRepositoryStandard{}
	)

	// import costs and export to single file for all of them
	path := filepath.Join(dataDir, "aws_costs")
	pattern := path + "/*.json"
	files, _ := filepath.Glob(pattern)

	for _, file := range files {
		old := []*v1.AwsCost{}
		fmt.Println(file)
		structs.UnmarshalFile(file, &old)

		for _, prior := range old {
			n := prior.V2()
			costs = append(costs, n)
		}
	}
	destination := filepath.Join(convertedDir, "aws_costs.json")
	structs.ToFile(costs, destination)

	// Import uptime
	path = filepath.Join(dataDir, "aws_uptime")
	pattern = path + "/*.json"
	files, _ = filepath.Glob(pattern)

	for _, file := range files {
		old := []*v1.AwsUptime{}
		fmt.Println(file)
		structs.UnmarshalFile(file, &old)

		for _, prior := range old {
			n := prior.V2()
			uptimes = append(uptimes, n)
		}
	}
	destination = filepath.Join(convertedDir, "aws_uptime.json")
	structs.ToFile(uptimes, destination)

	// Import standards
	path = filepath.Join(dataDir, "github_standards")
	pattern = path + "/*.json"
	files, _ = filepath.Glob(pattern)

	for _, file := range files {
		old := []*v1.GithubStandard{}
		fmt.Println(file)
		structs.UnmarshalFile(file, &old)

		for _, prior := range old {
			n := prior.V2()
			standards = append(standards, n)
		}
	}
	destination = filepath.Join(convertedDir, "github_standards.json")
	structs.ToFile(standards, destination)

}

func Run() (err error) {
	var (
		sess *session.Session
		svc  *s3.S3
	)

	if sess, err = awssession.New(awsCfg); err != nil {
		slog.Error("[convertor] aws session failed", slog.String("err", err.Error()))
		return
	}
	svc = s3.New(sess)
	Download(sess, svc)
	ConvertV1s()

	return

}

func main() {
	Run()
}

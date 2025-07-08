package existing

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
)

var t = true
var f = false
var assetid int64 = 2003
var ids = []int64{1000, 1001, 1002, 1003}
var nm = "test_accounts_v1.json"
var nms = []string{"R00", "R01", "R02", "R03"}
var mime = "application/json"
var rel = &github.RepositoryRelease{
	ID:         &ids[3],
	Name:       &nms[3],
	Draft:      &f,
	Prerelease: &t,
	Assets: []*github.ReleaseAsset{
		{ID: &assetid, Name: &nm, ContentType: &mime},
	},
}
var rels = []*github.RepositoryRelease{
	{ID: &ids[0], Name: &nms[0], Draft: &t, Prerelease: &t, Assets: []*github.ReleaseAsset{}},
	{ID: &ids[1], Name: &nms[1], Draft: &f, Prerelease: &f, Assets: []*github.ReleaseAsset{}},
	{ID: &ids[2], Name: &nms[2], Draft: &f, Prerelease: &t, Assets: []*github.ReleaseAsset{}},
	rel,
}

// mockedGitHubClient is a mock used for testing functions by returning fixed values
// so we dont have to makle api calls to github during testing
type mockedGitHubClient struct{}

func (self *mockedGitHubClient) ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) (all []*github.RepositoryRelease, resp *github.Response, err error) {
	all = rels
	resp = &github.Response{
		NextPage: 0,
	}
	return
}

func (self *mockedGitHubClient) GetLatestRelease(ctx context.Context, owner, repo string) (release *github.RepositoryRelease, resp *github.Response, err error) {
	release = rel
	return
}

// DownloadReleaseAsset is a dummy function that mimics downloading a specific id to a known file localtion
func (self *mockedGitHubClient) DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error) {
	var content = `[{
		"id": "001A",
		"name": "dev",
		"billing_unit": "TEAM-A",
		"label": "A",
		"environment": "development",
		"type": "aws",
		"uptime_tracking": true
	}]`
	rc = io.NopCloser(bytes.NewBuffer([]byte(content)))
	return
}

// mockedRepositoryS3BucketDownloader provides a mocked version of DownloadBucket that writes a dummy cost file to a
// known location and returns that as the file path
type mockedRepositoryS3BucketDownloader struct{}

// DownloadBucket generates a file with dummy cost data in to for testing inserts
func (self *mockedRepositoryS3BucketDownloader) DownloadBucket(client awsr.ClientS3ListAndGetter, bucket string, prefix string, directory string) (downloaded []string, err error) {
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

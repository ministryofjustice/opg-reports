package existing

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"opg-reports/report/internal/repository/awsr"
	"opg-reports/report/internal/utils"
	"os"
	"path/filepath"

	"github.com/google/go-github/v62/github"
)

var simpleRelID = utils.Ptr[int64](900000)
var simpleReleases = []*github.RepositoryRelease{
	{
		ID:         utils.Ptr[int64](90001),
		Name:       utils.Ptr("previous prerelease"),
		Draft:      utils.Ptr(false),
		Prerelease: utils.Ptr(true),
		TagName:    utils.Ptr("v1.0.1-prerelease"),
		Assets: []*github.ReleaseAsset{
			{
				ID:          utils.Ptr[int64](2001),
				Name:        utils.Ptr("asset.txt"),
				ContentType: utils.Ptr("text/plain"),
			},
			{
				ID:          utils.Ptr[int64](2002),
				Name:        utils.Ptr("asset.json"),
				ContentType: utils.Ptr("text/jsin"),
			},
		},
	},
	{
		ID:         simpleRelID,
		Name:       utils.Ptr("active release with assets"),
		Draft:      utils.Ptr(false),
		Prerelease: utils.Ptr(false),
		TagName:    utils.Ptr("v1.0.1"),
		Assets: []*github.ReleaseAsset{
			{
				ID:          utils.Ptr[int64](1001),
				Name:        utils.Ptr("asset.txt"),
				ContentType: utils.Ptr("text/plain"),
			},
			{
				ID:          utils.Ptr[int64](1002),
				Name:        utils.Ptr("asset.json"),
				ContentType: utils.Ptr("text/plain"),
			},
		},
	},
	{
		ID:         utils.Ptr[int64](89000),
		Name:       utils.Ptr("previous draft"),
		Draft:      utils.Ptr(true),
		Prerelease: utils.Ptr(false),
		TagName:    utils.Ptr("v1.0.1-draft"),
		Assets: []*github.ReleaseAsset{
			{
				ID:          utils.Ptr[int64](3001),
				Name:        utils.Ptr("asset.txt"),
				ContentType: utils.Ptr("text/plain"),
			},
			{
				ID:          utils.Ptr[int64](3002),
				Name:        utils.Ptr("asset.json"),
				ContentType: utils.Ptr("text/plain"),
			},
		},
	},
	{
		ID:         utils.Ptr[int64](80001),
		Name:       utils.Ptr("previous release"),
		Draft:      utils.Ptr(false),
		Prerelease: utils.Ptr(false),
		TagName:    utils.Ptr("v1.0.0"),
	},
}

type mockClientRepositoryReleaseListReleases struct{}

func (self *mockClientRepositoryReleaseListReleases) ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) (releases []*github.RepositoryRelease, resp *github.Response, err error) {
	releases = simpleReleases
	resp = &github.Response{NextPage: 0}
	return
}

func (self *mockClientRepositoryReleaseListReleases) DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error) {
	var asset *github.ReleaseAsset
	var content = `
[{
    "id": "500000067891",
    "name": "My production",
    "billing_unit": "Team A",
    "label": "prod",
    "environment": "production",
    "type": "aws",
    "uptime_tracking": true
	},
	{
    "id": "500000067891",
    "name": "My dev",
    "billing_unit": "Team A",
    "label": "dev",
    "environment": "development",
    "type": "aws",
    "uptime_tracking": true
}]
`
	for _, rel := range simpleReleases {
		for _, a := range rel.Assets {
			if *a.ID == id {
				asset = a
				break
			}
		}
	}

	if asset == nil {
		return
	}
	// content is name of the asset
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

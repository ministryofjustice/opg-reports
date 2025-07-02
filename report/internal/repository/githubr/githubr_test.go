package githubr

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v62/github"
)

// check interfaces are correct
var _ ReleaseRepository = &Repository{}

// mockedClient is a dummy stuck used for testing functions by using fixed values
// so we dont have to makle api calls to github during testing
type mockedClient struct{}

var t = true
var f = false
var id int64 = 2003
var ids = []int64{1000, 1001, 1002, 1003}
var nm = "test_accounts_v1.json"
var nms = []string{"R00", "R01", "R02", "R03"}
var rel = &github.RepositoryRelease{
	ID:         &ids[3],
	Name:       &nms[3],
	Draft:      &f,
	Prerelease: &t,
	Assets: []*github.ReleaseAsset{
		{ID: &id, Name: &nm},
	},
}
var rels = []*github.RepositoryRelease{
	{ID: &ids[0], Name: &nms[0], Draft: &t, Prerelease: &t, Assets: []*github.ReleaseAsset{}},
	{ID: &ids[1], Name: &nms[1], Draft: &f, Prerelease: &f, Assets: []*github.ReleaseAsset{}},
	{ID: &ids[2], Name: &nms[2], Draft: &f, Prerelease: &t, Assets: []*github.ReleaseAsset{}},
	rel,
}

func (self *mockedClient) ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) (all []*github.RepositoryRelease, resp *github.Response, err error) {
	all = rels
	resp = &github.Response{
		NextPage: 0,
	}
	return
}

func (self *mockedClient) GetLatestRelease(ctx context.Context, owner, repo string) (release *github.RepositoryRelease, resp *github.Response, err error) {
	release = rel
	return
}

// DownloadReleaseAsset is a dummy function that mimics downloading a specific id to a known file localtion
func (self *mockedClient) DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error) {
	var content = `[{
		"id": "001A",
		"name": "dev",
		"billing_unit": "TeamA",
		"label": "A",
		"environment": "development",
		"type": "aws",
		"uptime_tracking": true
	}]`
	rc = io.NopCloser(bytes.NewBuffer([]byte(content)))

	return
}

package githubr

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/utils"

	"github.com/google/go-github/v77/github"
)

// simple test release list
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
	rc = io.NopCloser(bytes.NewBuffer([]byte(*asset.Name)))

	return
}

func TestGithubrDownloadRepositoryReleaseAsset(t *testing.T) {

	var (
		err        error
		dir        = t.TempDir()
		ctx        = context.TODO()
		log        = utils.Logger("ERROR", "TEXT")
		conf       = config.NewConfig()
		owner      = "alphagov"
		repository = "govuk-frontend"
		opts       = &GetRepositoryReleaseOptions{
			ExcludePrereleases: true,
			ExcludeDraft:       true,
			ExcludeNoAssets:    true,
			ReleaseTag:         "v(.*)",
			UseRegex:           true,
		}
		aopts = &DownloadRepositoryReleaseAssetOptions{
			AssetName: "(.*)json",
			UseRegex:  true,
		}
		client = &mockClientRepositoryReleaseListReleases{}
		// client = DefaultClient(conf).Repositories
	)

	store := Default(ctx, log, conf)
	rel, err := store.GetRepositoryRelease(client, owner, repository, opts)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	// now try to download it
	_, localFile, err := store.DownloadRepositoryReleaseAsset(client, owner, repository, rel, dir, aopts)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if !utils.FileExists(localFile) {
		t.Errorf("downloaded file not found: %s", localFile)
	}
}

func TestGithubrGetRepositoryReleases(t *testing.T) {
	var (
		err        error
		ctx        = context.TODO()
		log        = utils.Logger("ERROR", "TEXT")
		conf       = config.NewConfig()
		owner      = "alphagov"
		repository = "govuk-frontend"
		opts       = &GetRepositoryReleaseOptions{
			ExcludePrereleases: true,
			ExcludeDraft:       true,
			ExcludeNoAssets:    true,
			ReleaseTag:         "v1*",
			UseRegex:           true,
		}
	)

	store := Default(ctx, log, conf)
	rels, err := store.GetRepositoryReleases(&mockClientRepositoryReleaseListReleases{}, owner, repository, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(rels) != len(simpleReleases) {
		t.Errorf("incorrect number returned")
	}
	// check the latest release is returned
	rels, err = store.GetRepositoryReleases(&mockClientRepositoryReleaseListReleases{}, owner, repository, opts)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if rels[0].ID != simpleRelID {
		t.Errorf("ordering failed")
	}

}

type optionValidationFixture struct {
	Release        *github.RepositoryRelease
	Options        *GetRepositoryReleaseOptions
	ExpectedResult bool
}

func TestGithubrGetRepositoryReleaseOptionsValidation(t *testing.T) {
	var inc bool
	var tests = []*optionValidationFixture{
		// Is DRAFT, but we're excluding those
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-01"), Draft: utils.Ptr(true), Prerelease: utils.Ptr(false)},
			Options:        &GetRepositoryReleaseOptions{ExcludePrereleases: true, ExcludeDraft: true},
			ExpectedResult: false,
		},
		// Is Prerelease, but excluding those
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-02"), Draft: utils.Ptr(false), Prerelease: utils.Ptr(true)},
			Options:        &GetRepositoryReleaseOptions{ExcludePrereleases: true, ExcludeDraft: true},
			ExpectedResult: false,
		},
		// Has no assets, so should be excluded
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-03"), Draft: utils.Ptr(false), Prerelease: utils.Ptr(false)},
			Options:        &GetRepositoryReleaseOptions{ExcludePrereleases: true, ExcludeDraft: true, ExcludeNoAssets: true},
			ExpectedResult: false,
		},
		// Is neither a draft or a prerelease, so should be returned
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-04"), Draft: utils.Ptr(false), Prerelease: utils.Ptr(false)},
			Options:        &GetRepositoryReleaseOptions{ExcludePrereleases: true, ExcludeDraft: true},
			ExpectedResult: true,
		},
		// Should be found by exact match on name
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-05"), Draft: utils.Ptr(false), Prerelease: utils.Ptr(false)},
			Options:        &GetRepositoryReleaseOptions{ReleaseName: "test-05"},
			ExpectedResult: true,
		},
		// Should not be found by exact match
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-06"), Draft: utils.Ptr(false), Prerelease: utils.Ptr(false)},
			Options:        &GetRepositoryReleaseOptions{ReleaseName: "test-0*"},
			ExpectedResult: false,
		},
		// Should be found regex pattern
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-07"), Draft: utils.Ptr(false), Prerelease: utils.Ptr(false)},
			Options:        &GetRepositoryReleaseOptions{ReleaseName: "test-0*", UseRegex: true},
			ExpectedResult: true,
		},
		// Should not be found - name matches, but no assets
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-08"), Draft: utils.Ptr(false), Prerelease: utils.Ptr(false)},
			Options:        &GetRepositoryReleaseOptions{ReleaseName: "test-0*", UseRegex: true, ExcludeNoAssets: true},
			ExpectedResult: false,
		},
		// Should not be found - name matches, but is draft
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-09"), Draft: utils.Ptr(true), Prerelease: utils.Ptr(false)},
			Options:        &GetRepositoryReleaseOptions{ReleaseName: "test-0*", UseRegex: true, ExcludeNoAssets: true, ExcludeDraft: true, ExcludePrereleases: true},
			ExpectedResult: false,
		},
		// Should not be found - name matches, but is draft
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-10"), Draft: utils.Ptr(true), Prerelease: utils.Ptr(false), Assets: []*github.ReleaseAsset{{ID: utils.Ptr[int64](100), Name: utils.Ptr("meta.tar.bz")}}},
			Options:        &GetRepositoryReleaseOptions{ReleaseName: "test-0*", UseRegex: true, ExcludeNoAssets: true, ExcludeDraft: true, ExcludePrereleases: true},
			ExpectedResult: false,
		},
		// Should be found - name matches, has assets
		{
			Release:        &github.RepositoryRelease{Name: utils.Ptr("test-11"), Draft: utils.Ptr(false), Prerelease: utils.Ptr(false), Assets: []*github.ReleaseAsset{{ID: utils.Ptr[int64](100), Name: utils.Ptr("meta.tar.bz")}}},
			Options:        &GetRepositoryReleaseOptions{ReleaseName: "test-0*", UseRegex: true, ExcludeNoAssets: true, ExcludeDraft: true, ExcludePrereleases: true},
			ExpectedResult: true,
		},
	}
	// is draft, should be skipped
	for _, test := range tests {
		inc = repositoryReleaseMeetsCriteria(test.Release, test.Options)
		if inc != test.ExpectedResult {
			t.Errorf("expected [%v], actual [%v] see test [%v]", test.ExpectedResult, inc, *test.Release.Name)
		}
	}
}

package githubr

import (
	"context"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/utils"

	"github.com/google/go-github/v62/github"
)

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
		inc = test.Options.Validate(test.Release)
		if inc != test.ExpectedResult {
			t.Errorf("expected [%v], actual [%v] see test [%v]", test.ExpectedResult, inc, *test.Release.Name)
		}
	}
}

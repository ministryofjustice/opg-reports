package githubr

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"opg-reports/report/config"
	"opg-reports/report/internal/utils"
	"testing"

	"github.com/google/go-github/v77/github"
)

// mock the client for testing
type mockClientOwnership struct{}

func (self *mockClientOwnership) DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (rc io.ReadCloser, resp *github.Response, err error) {
	resp = &github.Response{NextPage: 0, Response: &http.Response{StatusCode: http.StatusOK}}
	// resp.StatusCode = http.StatusOK
	content := `* @ministryofjustice/a @ministryofjustice/b
/.github/  @ministryofjustice/a
	`
	rc = io.NopCloser(bytes.NewBuffer([]byte(content)))
	return
}
func (self *mockClientOwnership) ListTeams(ctx context.Context, owner, repo string, opts *github.ListOptions) (teams []*github.Team, resp *github.Response, err error) {
	resp = &github.Response{NextPage: 0, Response: &http.Response{StatusCode: http.StatusOK}}
	resp.StatusCode = http.StatusOK
	teams = []*github.Team{
		{
			Name: utils.Ptr("A"),
			Slug: utils.Ptr("a"),
			Parent: &github.Team{
				Name: utils.Ptr("opg"),
				Slug: utils.Ptr("opg"),
			},
			Organization: &github.Organization{
				Login: utils.Ptr("ministryofjustice"),
			},
		},
		{
			Name: utils.Ptr("B"),
			Slug: utils.Ptr("b"),
			Parent: &github.Team{
				Name: utils.Ptr("opg"),
				Slug: utils.Ptr("opg"),
			},
			Organization: &github.Organization{
				Login: utils.Ptr("ministryofjustice"),
			},
		},
		{
			Name: utils.Ptr("C"),
			Slug: utils.Ptr("c"),
			Organization: &github.Organization{
				Login: utils.Ptr("ministryofjustice"),
			},
			Parent: &github.Team{
				Name: utils.Ptr("laa"),
				Slug: utils.Ptr("laa"),
			},
		},
	}
	return
}

func TestGithubrGetRepositoryOwners(t *testing.T) {
	var (
		err  error
		res  []string
		ctx  = context.TODO()
		log  = utils.Logger("ERROR", "TEXT")
		conf = config.NewConfig()
		rp   = Default(ctx, log, conf)
		// client = DefaultClient(conf).Repositories
		client = &mockClientOwnership{}
		repo   = &github.Repository{
			Name:     utils.Ptr("opg-modernising-lpa"),
			FullName: utils.Ptr("ministryofjustice/opg-modernising-lpa"),
			Owner: &github.User{
				Login: utils.Ptr("ministryofjustice"),
			},
		}
	)
	res, err = rp.GetRepositoryOwners(client, repo, &GetRepositoryOwnerOptions{})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 3 {
		t.Errorf("not all owners returned")
	}
	// test the exclusion logic
	res, err = rp.GetRepositoryOwners(client, repo, &GetRepositoryOwnerOptions{Exclude: []string{"ministryofjustice/a"}})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 2 {
		t.Errorf("incorrect owners returned, exclusion failed")
	}

	// check filter only returns opg parents and not `c` dummy
	res, err = rp.GetRepositoryOwners(client, repo, &GetRepositoryOwnerOptions{FilterByParent: "opg"})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 2 {
		t.Errorf("filtering of owners failed")
	}

	// check filter and the exclusion
	res, err = rp.GetRepositoryOwners(client, repo, &GetRepositoryOwnerOptions{FilterByParent: "opg", Exclude: []string{"ministryofjustice/a"}})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 1 {
		t.Errorf("filtering of owners with exclusions failed")
	}

}

func TestGithubrGetTeamsForRepository(t *testing.T) {

	var (
		err  error
		res  []*github.Team
		ctx  = context.TODO()
		log  = utils.Logger("ERROR", "TEXT")
		conf = config.NewConfig()
		rp   = Default(ctx, log, conf)
		// client = DefaultClient(conf).Repositories
		client = &mockClientOwnership{}
		repo   = &github.Repository{
			Name:     utils.Ptr("opg-use-an-lpa"),
			FullName: utils.Ptr("ministryofjustice/opg-use-an-lpa"),
			Owner: &github.User{
				Login: utils.Ptr("ministryofjustice"),
			},
		}
	)
	// check no filter
	res, err = rp.GetTeamsForRepository(client, repo, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 3 {
		t.Errorf("unexpected length returned, filter issue")
	}

	// check filter by parent works
	res, err = rp.GetTeamsForRepository(client, repo, &GetRepositoryOwnerOptions{FilterByParent: "opg"})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 2 {
		t.Errorf("unexpected length returned, filter failed")
	}

	// check exclusion works
	res, err = rp.GetTeamsForRepository(client, repo, &GetRepositoryOwnerOptions{Exclude: []string{"ministryofjustice/a"}})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 2 {
		t.Errorf("unexpected length returned, filter failed")
	}

}

func TestGithubrGetCodeOwnersForRepository(t *testing.T) {
	var (
		err  error
		res  []string
		ctx  = context.TODO()
		log  = utils.Logger("ERROR", "TEXT")
		conf = config.NewConfig()
		rp   = Default(ctx, log, conf)
		// client = DefaultClient(conf).Repositories
		client = &mockClientOwnership{}
		repo   = &github.Repository{
			Name:     utils.Ptr("opg-modernising-lpa"),
			FullName: utils.Ptr("ministryofjustice/opg-modernising-lpa"),
			Owner: &github.User{
				Login: utils.Ptr("ministryofjustice"),
			},
		}
	)

	res, err = rp.GetCodeOwnersForRepository(client, repo, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 2 {
		t.Errorf("incorrect length from codeowners")
	}

	res, err = rp.GetCodeOwnersForRepository(client, repo, &GetRepositoryOwnerOptions{Exclude: []string{"ministryofjustice/b"}})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 1 {
		t.Errorf("incorrect length from codeowners, exclude failed")
	}
}

type testCodeOwnerLines struct {
	Lines    []string
	Expected int
}

func TestGithubrCodeOwnersFromLines(t *testing.T) {

	var tests = []*testCodeOwnerLines{
		{
			Expected: 2,
			Lines: []string{
				"* @ministryofjustice/opg @ministryofjustice/opg-webops",
				"/.github/  @ministryofjustice/opg-webops",
			},
		},
		{
			Expected: 1,
			Lines: []string{
				"* @ministryofjustice/opg-webops",
			},
		},
	}

	for _, test := range tests {
		res := codeOwnersFromLines(test.Lines)
		if test.Expected != len(res) {
			t.Errorf("unexpected number of codeowners returned expected [%v] actual [%v]", test.Expected, len(res))
		}
	}

}

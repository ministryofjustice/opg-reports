package ghteamrepositories

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/app/internal/convert"
	"opg-reports/app/internal/ghdata/ghclient"
	"opg-reports/app/internal/ghdata/ghfilters"
	"testing"

	"github.com/google/go-github/v87/github"
)

// check mock works for the client
var _ Client = &mockTeamService{}

// mockTeamService client to fetch a fixed set of repositories which can then be tested with filters
type mockTeamService struct{}

func (self *mockTeamService) ListTeamReposBySlug(ctx context.Context, orgSlug string, teamSlug string, opts *github.ListOptions) (repos []*github.Repository, response *github.Response, err error) {

	response = &github.Response{
		Response: &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
		},
		NextPage: 0,
	}

	owner := &github.User{
		Login: convert.Ptr("ministryofjustice-test-a"),
		Name:  convert.Ptr("ministryofjustice-test-a"),
	}

	repos = []*github.Repository{
		// public repo thats active
		{
			FullName:      convert.Ptr(fmt.Sprintf("%s/%s", *owner.Login, "opg-example-public-active")),
			Name:          convert.Ptr("opg-example-public-active"),
			DefaultBranch: convert.Ptr("main"),
			Archived:      convert.Ptr(false),
			Visibility:    convert.Ptr("public"),
			Private:       convert.Ptr(false),
			Owner:         owner,
		},
		// private repo thats active
		{
			FullName:      convert.Ptr(fmt.Sprintf("%s/%s", *owner.Login, "opg-example-private-active")),
			Name:          convert.Ptr("opg-example-private-active"),
			DefaultBranch: convert.Ptr("main"),
			Archived:      convert.Ptr(false),
			Visibility:    convert.Ptr("private"),
			Private:       convert.Ptr(true),
			Owner:         owner,
		},
		// archived public repo
		{
			FullName:      convert.Ptr(fmt.Sprintf("%s/%s", *owner.Login, "opg-example-public-archived")),
			Name:          convert.Ptr("opg-example-public-archived"),
			DefaultBranch: convert.Ptr("main"),
			Archived:      convert.Ptr(true),
			Visibility:    convert.Ptr("public"),
			Private:       convert.Ptr(false),
			Owner:         owner,
		},
	}

	return
}

// TestGHTeamRepositoriesGetDataMocked uses mocked client to provide
// preset data
func TestGHTeamRepositoriesGetDataMocked(t *testing.T) {
	var (
		err     error
		client  *mockTeamService
		res     []*github.Repository
		skipped []any
		src     *Source[*mockTeamService, *github.Repository]
		ctx     context.Context = t.Context()
		cfg     *Config         = &Config{
			OrganisationSlug: "ministryofjustice-test",
			TeamSlug:         "opg",
		}
	)

	client = &mockTeamService{}
	// now create the data source
	src, err = New[*mockTeamService, *github.Repository](ctx, client, cfg)
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}
	// test fetching the data
	res, _, err = src.GetData()
	if len(res) <= 0 {
		t.Errorf("failed to find any repositories.")
	}

	// now test this with a filter removes archived values
	src, err = New[*mockTeamService, *github.Repository](ctx, client, cfg, &ghfilters.ExcludeArchivedRepository{})
	// test fetching the data
	res, skipped, err = src.GetData()
	if len(res) <= 0 {
		t.Errorf("failed to find any repositories.")
	}
	// should have some repos skipped
	if len(skipped) == 0 {
		t.Errorf("expected some repositories to be skipped as they are archived..")
	}

	// now check that results are not archived
	for _, r := range res {
		if *r.Archived {
			t.Errorf("unexpected archived repo: [%s]", *r.FullName)
		}
	}
}

// TestGHTeamRepositoriesGetDataActual uses real api connection and client
// to fetch data
func TestGHTeamRepositoriesGetDataActual(t *testing.T) {
	var (
		err     error
		client  *github.Client
		res     []*github.Repository
		skipped []any
		src     *Source[*github.TeamsService, *github.Repository]
		token   string          = ghclient.Token()
		ctx     context.Context = t.Context()
		cfg     *Config         = &Config{
			OrganisationSlug: "ministryofjustice",
			TeamSlug:         "opg",
		}
	)
	// if theres no github token, skip this test
	if token == "" {
		t.SkipNow()
	}
	// create the client
	client, err = ghclient.New(ctx, token)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	// now create the data source
	src, err = New[*github.TeamsService, *github.Repository](ctx, client.Teams, cfg)
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}
	// test fetching the data
	res, _, err = src.GetData()
	if len(res) <= 0 {
		t.Errorf("failed to find any repositories.")
	}

	// now test this with a filter removes archived values
	src, err = New[*github.TeamsService, *github.Repository](ctx, client.Teams, cfg, &ghfilters.ExcludeArchivedRepository{})
	// test fetching the data
	res, skipped, err = src.GetData()
	if len(res) <= 0 {
		t.Errorf("failed to find any repositories.")
	}
	// should have some repos skipped
	if len(skipped) == 0 {
		t.Errorf("expected some repositories to be skipped as they are archived..")
	}

	// now check that results are not archived
	for _, r := range res {
		if *r.Archived {
			t.Errorf("unexpected archived repo: [%s]", *r.FullName)
		}
	}
}

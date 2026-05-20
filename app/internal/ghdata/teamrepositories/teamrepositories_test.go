package teamrepositories

import (
	"context"
	"opg-reports/app/internal/ghdata/ghclient"
	"opg-reports/app/internal/ghdata/ghfilters"
	"testing"

	"github.com/google/go-github/v87/github"
)

// TestTeamRepositoriesGetDataActual uses real api connection and client
// to fetch data
func TestTeamRepositoriesGetDataActual(t *testing.T) {
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

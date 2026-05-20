package teamrepositories

import (
	"context"
	"opg-reports/app/internal/fmtx"
	"opg-reports/app/internal/ghdata/ghclient"
	"opg-reports/app/internal/ghdata/ghconfig"
	"testing"

	"github.com/google/go-github/v87/github"
)

// TestGetDataActual uses real api connection and client
// to fetch repositories
func TestGetDataActual(t *testing.T) {
	var (
		err    error
		client *github.Client
		res    []*github.Repository
		src    *Source[*github.TeamsService, *github.Repository]
		token  string           = ghclient.Token()
		ctx    context.Context  = t.Context()
		cfg    *ghconfig.Config = ghconfig.New("ministryofjustice", "opg")
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

	res, err = src.GetData()
	fmtx.Printj(res)

	t.FailNow()
}

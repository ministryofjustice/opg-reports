package codebase

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/ptr"
	"os"
	"testing"

	"github.com/google/go-github/v81/github"
)

// mockGetter is a mocked client that returns successful results
type mockGetter struct{}

func (self *mockGetter) ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) (repos []*github.Repository, resp *github.Response, err error) {

	repos = []*github.Repository{
		{Name: ptr.Ptr("mock-repo-a"), FullName: ptr.Ptr("mock/mock-repo-a"), HTMLURL: ptr.Ptr("https://test.local/mock/mock-repo-a"), Archived: ptr.Ptr(false)},
		{Name: ptr.Ptr("mock-repo-b"), FullName: ptr.Ptr("mock/mock-repo-b"), HTMLURL: ptr.Ptr("https://test.local/mock/mock-repo-b"), Archived: ptr.Ptr(false)},
		{Name: ptr.Ptr("mock-repo-c"), FullName: ptr.Ptr("mock/mock-repo-c"), HTMLURL: ptr.Ptr("https://test.local/mock/mock-repo-c"), Archived: ptr.Ptr(true)},
	}
	resp = &github.Response{
		NextPage: 0,
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
	}
	return
}

func TestDomainCodebasesWithMock(t *testing.T) {

	var (
		err    error
		client *mockGetter     = &mockGetter{}
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		opts   *Options        = &Options{ExcludeArchived: true, OrgSlug: "mock", ParentTeam: "test"}
		data   []*codebasemodels.Codebase
	)

	data, err = GetCodebases(ctx, log, client, opts)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
	if len(data) != 2 {
		t.Errorf("expected exactly two teams in the list.")
	}

}

func TestDomainCodebasesWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		data   []*codebasemodels.Codebase
	)

	if os.Getenv("GH_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GH_TOKEN"))
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		opts := &Options{ExcludeArchived: true, OrgSlug: "ministryofjustice", ParentTeam: "opg"}
		data, err = GetCodebases(ctx, log, client.Teams, opts)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		if len(data) < 5 {
			t.Errorf("expected more teams in the list")
		}

	} else {
		t.SkipNow()
	}
}

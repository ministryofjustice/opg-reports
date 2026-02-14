package codeownergetter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/ptr"
	"os"
	"testing"

	"github.com/google/go-github/v81/github"
)

type mockGetter struct{}

func (self *mockGetter) ListTeams(ctx context.Context, owner, repo string, opts *github.ListOptions) (teams []*github.Team, resp *github.Response, err error) {

	parent := &github.Team{
		Slug: ptr.Ptr("opg"),
	}

	teams = []*github.Team{
		{
			Name:    ptr.Ptr("mock-other"),
			Slug:    ptr.Ptr("mock-other"),
			HTMLURL: ptr.Ptr("https://github.com/orgs/" + owner + "/mock-other"),
		},
		{
			Name:    ptr.Ptr("mock-team-a"),
			Slug:    ptr.Ptr("mock-team-a"),
			HTMLURL: ptr.Ptr("https://github.com/orgs/" + owner + "/mock-team-a"),
			Parent:  parent,
		},
		{
			Name:    ptr.Ptr("mock-team-b"),
			Slug:    ptr.Ptr("mock-team-b"),
			HTMLURL: ptr.Ptr("https://github.com/orgs/" + owner + "/mock-team-b"),
			Parent:  parent,
		},
	}
	// add extra team for mock-repo-b
	if repo == "mock-repo-b" {
		teams = append(teams, &github.Team{
			Name:    ptr.Ptr("mock-team-c"),
			Slug:    ptr.Ptr("mock-team-c"),
			HTMLURL: ptr.Ptr("https://github.com/orgs/" + owner + "/mock-team-c"),
			Parent:  parent,
		})
	}

	resp = &github.Response{
		NextPage: 0,
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
	}
	return
}

func (self *mockGetter) DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (buff io.ReadCloser, resp *github.Response, err error) {

	content := fmt.Sprintf(`* @%s/mock-team-a @%s/mock-team-b%s`, owner, owner, "\n")
	if repo == "mock-repo-b" {
		content += fmt.Sprintf(`.github/  @%s/mock-team-codeowner%s`, owner, "\n")
	}
	buff = io.NopCloser(bytes.NewBuffer([]byte(content)))

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
		opts   *Input
		org    string                       = "ministryofjustice"
		client *mockGetter                  = &mockGetter{}
		ctx    context.Context              = t.Context()
		log    *slog.Logger                 = logger.New("error")
		data   []*codeownermodels.Codeowner = []*codeownermodels.Codeowner{}
	)

	opts = &Input{
		OrgSlug:    org,
		ParentTeam: "mock-parent",
		Codebases: []*codebasemodels.Codebase{
			{FullName: org + "/mock-repo-a", Name: "mock-repo-a"},
			{FullName: org + "/mock-repo-b", Name: "mock-repo-b"},
			{FullName: org + "/mock-repo-c", Name: "mock-repo-c"},
		},
	}

	data, err = GetCodeowners(ctx, log, client, opts)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
	if len(data) < 1 {
		t.Errorf("unexpected count of results")
	}

	foundOther := false // this team did not have correct parent..
	for _, row := range data {
		if row.Name == "mock-other" {
			foundOther = true
		}
	}
	if foundOther {
		t.Errorf("returned data included a mock team that should have been excluded.")
	}

	foundExtraBTeam := false
	for _, row := range data {
		if row.Name == "mock-team-c" && row.CodebaseFullName == org+"/mock-repo-b" {
			foundExtraBTeam = true
		}
	}
	if foundExtraBTeam {
		t.Errorf("returned data did not include extrea team for the b-repo.")
	}

}

func TestDomainCodeownerWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		opts   *Input
		ctx    context.Context              = t.Context()
		log    *slog.Logger                 = logger.New("error")
		result []*codeownermodels.Codeowner = []*codeownermodels.Codeowner{}
	)
	opts = &Input{
		OrgSlug:    "ministryofjustice",
		ParentTeam: "opg",
		Codebases: []*codebasemodels.Codebase{
			{FullName: "ministryofjustice/opg-lpa", Name: "opg-lpa"},
			{FullName: "ministryofjustice/opg-use-an-lpa", Name: "opg-use-an-lpa"},
			{FullName: "ministryofjustice/opg-data-lpa-store", Name: "opg-data-lpa-store"},
		},
	}
	if os.Getenv("GH_TOKEN") == "" {
		t.SkipNow()
	}

	client, err = ghclients.New(ctx, log, os.Getenv("GH_TOKEN"))
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

	result, err = GetCodeowners(ctx, log, client.Repositories, opts)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
	if len(result) <= 0 {
		t.Errorf("unexpected count of results")
	}

}

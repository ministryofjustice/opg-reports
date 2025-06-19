package gh

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/report/config"
)

type Repository struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

// connection is an internal helper to handle creating the client
func (self *Repository) connection() (client *github.Client, err error) {
	// handle empty configs
	if self.conf.Github == nil || self.conf.Github.Token == "" {
		return nil, fmt.Errorf("no github access token found in the config")
	}
	// get a rate limted version of the client
	limited, err := github_ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		return
	}
	client = github.NewClient(limited).WithAuthToken(self.conf.Github.Token)

	return
}

type ReleaseOptions struct {
	ExcludePrereleases bool
	ExcludeDraft       bool
	ExcludeNoAssets    bool
}

// getAllReleases returns all releases for a repository without any filtering
// - repositoryName should not include the organsiation name
func (self *Repository) getAllReleases(organisation string, repositoryName string) (releases []*github.RepositoryRelease, err error) {
	var (
		client *github.Client
		page   int                 = 1
		opts   *github.ListOptions = &github.ListOptions{PerPage: 200}
		log                        = self.log.With("repositoryName", repositoryName, "operation", "getAllReleases")
	)
	// get api client
	client, err = self.connection()
	if err != nil {
		return
	}
	// loop around the pagination
	for page > 0 {
		var response *github.Response
		var list []*github.RepositoryRelease
		// set the page number
		opts.Page = page
		// get all releases for the repository
		log.With("page", page).Debug("getting next page of releases")
		list, response, err = client.Repositories.ListReleases(self.ctx, organisation, repositoryName, opts)
		if err != nil {
			return
		}
		// if there items in the list, them merge into all
		if len(list) > 0 {
			releases = append(releases, list...)
		}
		// move to next page
		page = response.NextPage
	}

	return
}

// GetReleases returns all releases for a repository with some basic filtering options available.
//
// If options is nil (or all values are false) then all releases are returned.
func (self *Repository) GetReleases(organisation string, repositoryName string, options *ReleaseOptions) (releases []*github.RepositoryRelease, err error) {
	// setup log to be for this operation
	var log = self.log.With("repositoryName", repositoryName, "operation", "Releases")
	releases = []*github.RepositoryRelease{}
	// first, get all releases
	all, err := self.getAllReleases(organisation, repositoryName)

	// if there are no filter, then return everything
	if options == nil || (!options.ExcludeDraft && !options.ExcludeNoAssets && !options.ExcludePrereleases) {
		log.Debug("no filtering set, returning all releases directly")
		releases = all
		return
	}
	// there are filters, so look what to return
	log.Debug("filtering required, checking release values")
	// find only the required releases
	for _, release := range all {
		var keep bool = true
		// add release info to the log output
		var lg = log.With("draft", *release.Draft, "prerelease", *release.Prerelease, "assets", len(release.Assets), "id", *release.ID)
		// if this is a draft, but we are excluding them, swap keep to false
		if *release.Draft == true && options.ExcludeDraft == true {
			keep = false
		}
		// if this is a prerelease and we are excluding them, keep to false
		if *release.Prerelease == true && options.ExcludePrereleases == true {
			keep = false
		}
		// if there are no assets and we are excluding those without, the swap
		if len(release.Assets) == 0 && options.ExcludeNoAssets == true {
			keep = false
		}
		if keep {
			releases = append(releases, release)
		}
		lg.With("keep", keep).Debug("release checked")
	}

	return
}

// GetLatestRelease returns the latest published release for a repository.
// If you are looking for a prerelease / draft then use GetReleases with options configured
// - repositoryName should not include the organsiation name
func (self *Repository) GetLatestRelease(organisation string, repositoryName string) (release *github.RepositoryRelease, err error) {
	var client *github.Client
	var log = self.log.With("repositoryName", repositoryName, "operation", "GetLatestRelease")
	// get api client
	client, err = self.connection()
	if err != nil {
		return
	}

	// get just 1 release - this should be the latest
	log.Debug("getting last release")
	release, _, err = client.Repositories.GetLatestRelease(self.ctx, organisation, repositoryName)

	return
}

func New(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository, err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for github repository")
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for github repository")
	}
	log = log.WithGroup("github")
	rp = &Repository{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}

	return
}

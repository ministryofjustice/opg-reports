package repos

import (
	"context"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/times"
	"time"

	"github.com/google/go-github/v84/github"
)

// GetMergedPRs fetches all pull requests for a repository.
func GetMergedPRs(ctx context.Context, client prClient, repo *github.Repository, in *Args) (prs []*github.PullRequest, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "repos", "func", "GetMergedPRs", "repo", *repo.Name)

	log.Debug("getting merged pull requests ...")
	prs, err = paginatedMergedPRs(ctx, client, repo, in, &github.PullRequestListOptions{
		State:     "closed",
		Sort:      "created",
		Direction: "desc",
		Base:      *repo.DefaultBranch,
	})

	log.With("count", len(prs)).Debug("complete.")
	return
}

// paginatedMergedPRs iterrates over github results and merges results
// avoiding duplicates by checking the  run id
func paginatedMergedPRs(ctx context.Context, client prClient, repo *github.Repository, in *Args, opts *github.PullRequestListOptions) (result []*github.PullRequest, err error) {
	var (
		maxRetry int                           = 3
		page     int                           = 1
		all      map[int64]*github.PullRequest = map[int64]*github.PullRequest{}
		log      *slog.Logger                  = cntxt.GetLogger(ctx).With("package", "repos", "func", "paginatedMergedPRs")
	)
	result = []*github.PullRequest{}
	// force the max per page
	opts.PerPage = 100
	log.Debug("starting ....")
	for page > 0 {
		var (
			prs      []*github.PullRequest
			response *github.Response
			found    map[int64]*github.PullRequest
			retry    = 0
		)

		opts.Page = page
		log.Debug("getting page of results ...", "page", page)
		prs, response, err = client.List(ctx, *repo.Owner.Login, *repo.Name, opts)
		// simple re-try loop as we get sporadic failures
		for err != nil && retry < maxRetry {
			retry += 1
			log.Warn("error getting pull request data, retrying in 1 second ...", "err", err.Error())
			time.Sleep(time.Second * 1)
			prs, response, err = client.List(ctx, *repo.Owner.Login, *repo.Name, opts)
		}

		if err != nil {
			log.Error("error getting pull request data", "err", err.Error())
			return
		}
		// process all the prs and get next page
		found, page, err = processPRs(ctx, prs, in, response)
		if err != nil {
			return
		}
		// merge into main set
		for k, v := range found {
			all[k] = v
		}

	}
	// push from map to slice
	for _, pr := range all {
		result = append(result, pr)
	}
	log.With("count", len(result)).Debug("complete.")
	return
}

// processPRs
// works out next page
func processPRs(ctx context.Context, prs []*github.PullRequest, in *Args, response *github.Response) (found map[int64]*github.PullRequest, nextPage int, err error) {
	var (
		dateStart time.Time    = times.Add(in.DateStart, -1, times.SECOND)
		dateEnd   time.Time    = times.Add(in.DateEnd, 1, times.SECOND)
		log       *slog.Logger = cntxt.GetLogger(ctx).With("package", "repos", "func", "processPRs")
	)

	nextPage = response.NextPage
	found = map[int64]*github.PullRequest{}

	// all runs should have unique id
	for _, pr := range prs {
		var when time.Time
		// if this pr is not merged, then skip it
		if pr.MergeCommitSHA == nil || pr.MergedAt == nil || len(*pr.MergeCommitSHA) <= 0 {
			log.Debug("merge commit data missing, skipping pr")
			continue
		}
		// used created time for date range checking, as this cannot change
		when = pr.CreatedAt.Time
		// if its with the date range we want, add to the list
		if when.After(dateStart) && when.Before(dateEnd) {
			log.Debug("pr is with date range provided")
			found[*pr.ID] = pr
		}
		// if this was created before the date range, then break the paginatin fetching
		// by returning a next page of 0
		if when.Before(dateStart) {
			log.Debug("pr is outside of date range ...")
			nextPage = 0
			return
		}
	}
	return
}

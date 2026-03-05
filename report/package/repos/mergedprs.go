package repos

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/times"
	"time"

	"github.com/google/go-github/v84/github"
)

// GetMergedPRs fetches prs for the reporsitory between the dates stipulated via
// `DateStart` & `DateEnd` that target the default branch, are closed and merged.
//
// # Used as a back up for workflows as some codebases run in jenkins etc.
//
// Date restrictions are done via custom filtering as the api has no date params
func GetMergedPRs(ctx context.Context, client prClient, repo *github.Repository, in *Args) (prs []*github.PullRequest, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "repos", "func", "GetMergedPRs", "repo", *repo.Name)

	log.Debug("getting merged pull requests ...")
	prs, err = paginatedMergedPRs(ctx, client, repo, in, &github.PullRequestListOptions{
		State:     "closed",
		Sort:      "updated",
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
		var prs []*github.PullRequest
		var response *github.Response
		var found map[int64]*github.PullRequest
		var retry = 0
		log.Debug("getting page of results ...", "page", page)

		opts.Page = page
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
		// process the prs
		found, err = processPRs(ctx, prs, in)
		// if the prs are outside of date range, its not a real error, just break the fetch loop
		// otherwise return any other errors
		if err != nil && errors.Is(err, ErrPROutOfDateRange) {
			page = 0
			err = nil
		} else if err != nil {
			return
		} else {
			page = response.NextPage
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

func processPRs(ctx context.Context, prs []*github.PullRequest, in *Args) (found map[int64]*github.PullRequest, err error) {
	var dateStart time.Time = times.Add(in.DateStart, -1, times.SECOND)
	var dateEnd time.Time = times.Add(in.DateEnd, 1, times.SECOND)
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "repos", "func", "processPRs")

	found = map[int64]*github.PullRequest{}
	// all runs should have unique id
	for _, pr := range prs {
		var when time.Time = pr.UpdatedAt.Time
		// prs are sorted latest first, so if we have found an old one,
		// we break the loop and stop calling the api
		if pr.UpdatedAt.Time.Before(dateStart) {
			err = ErrPROutOfDateRange
			log.Debug("found pr older than start date; breaking loop ...")
			return
		}
		// if this pr is not merged, then skip it
		if pr.MergeCommitSHA == nil || pr.MergedAt == nil || len(*pr.MergeCommitSHA) <= 0 {
			log.Debug("merge commit data missing, skipping pr")
			continue
		}
		// as the List api call has no data filtering, we need to add our own...
		if when.After(dateStart) && when.Before(dateEnd) {
			log.Debug("pr is with date range provided")
			found[*pr.ID] = pr
		}
	}
	return
}

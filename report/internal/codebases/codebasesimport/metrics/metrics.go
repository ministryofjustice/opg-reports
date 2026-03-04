// Package metrics calls the github api to fetch data about time related events.
//
// This is very heavy on api calls and might hit rate limiter due to the number
// calls need to make to fetch workflow run data and pull requests
package metrics

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/codebases/codebasesimport/args"
	"opg-reports/report/internal/codebases/codebasesimport/clients"
	"opg-reports/report/package/cntxt"

	"github.com/google/go-github/v84/github"
)

// Raw stats entry
const InsertMetricsStatement string = `
INSERT INTO codebase_metrics (
	codebase,
	month,
	releases,
	releases_securityish,
	pr_count,
	pr_count_securityish,
	pr_stale_count,
	average_time_live,
	average_time_pr
) VALUES (
	:codebase,
	:month,
	:releases,
	:releases_securityish,
	:pr_count,
	:pr_count_securityish,
	:pr_stale_count,
	:average_time_live,
	:average_time_pr
)
ON CONFLICT (codebase,month) DO UPDATE SET
	releases=excluded.releases,
	releases_securityish=excluded.releases_securityish,
	pr_count=excluded.pr_count,
	pr_count_securityish=excluded.pr_count_securityish,
	pr_stale_count=excluded.pr_stale_count,
	average_time_live=excluded.average_time_live,
	average_time_live=excluded.average_time_live
RETURNING id
;
`

type Clients struct {
	Actions clients.ActionClient
	PR      clients.PRClient
}

type CodebaseMetric struct {
	Codebase string `json:"codebase,omitempty"` // full name of codebase
	Month    string `json:"month,omitempty"`    // month as YYYY-MM string

	Releases            int    `json:"releases,omitempty"`             // count of releases for this month
	ReleasesSecurityish int    `json:"releases_securityish,omitempty"` // count of releases for this month that seem to be security related
	AverageTimeLive     string `json:"average_time_live"`              // average time path to live workflow took (in milliseconds)

	PRCount            int    `json:"pr_count,omitempty"`   // count of all pull requests for the month
	PRCountSecurityish int    `json:"pr_count_securityish"` // count of all pr's that roughly relate to security (bots / keywords)
	PRStaleCount       int    `json:"pr_count_stale"`       // count of stale pull requests - open for longer than x days
	AverageTimePR      string `json:"average_time_pr"`      // average time a pull request workflow took (in milliseconds)
}

var ErrNoWorkflows = errors.New("no workflow runs for this repository.")

func HandleCodebaseMetrics(ctx context.Context, client *Clients, repositories []*github.Repository, in *args.Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "HandleCodebaseMetrics")
	var data []*CodebaseMetric = []*CodebaseMetric{}
	log.With("count", len(repositories)).Info("starting codebase metrics import ...")

	for _, repo := range repositories {
		log.Info("getting metric data for repository ...", "repo", *repo.Name)
		fmt.Println(*repo.Name)
		workflowMetrics(ctx, client.Actions, repo, in)
	}

	log.With("count", len(data)).Info("complete.")
	return
}

func emptyMetric(codebase string, month string) *CodebaseMetric {
	return &CodebaseMetric{
		Codebase:            codebase,
		Month:               month,
		Releases:            0,
		ReleasesSecurityish: 0,
		AverageTimeLive:     "0.0",
		PRCount:             0,
		PRCountSecurityish:  0,
		PRStaleCount:        0,
		AverageTimePR:       "0.0",
	}
}

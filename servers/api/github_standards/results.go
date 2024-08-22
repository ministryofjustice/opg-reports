package github_standards

import (
	"context"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

// getResults handles determining which query to call based on the get param values
func getResults(ctx context.Context, queries *ghs.Queries, archived string, team string) (results []ghs.GithubStandard, err error) {

	var teamF = ""
	var archivedF = ""
	results = []ghs.GithubStandard{}
	// -- fetch the get parameter values
	// team query, add the like logic here
	if team != "" {
		teamF = "%#" + team + "#%"
	}
	// archive query
	if archived != "" {
		archivedF = archived
	}
	// -- run queries
	if teamF != "" && archivedF != "" {
		// if both team and archive are set, use joined query
		results, err = queries.FilterByIsArchivedAndTeam(ctx, ghs.FilterByIsArchivedAndTeamParams{
			IsArchived: convert.BoolStringToInt(archivedF), Teams: teamF,
		})
	} else if archivedF != "" {
		// run for just archived - this is defaulted to 1
		results, err = queries.FilterByIsArchived(ctx, convert.BoolStringToInt(archivedF))
	} else if teamF != "" {
		// if only team is set, then return team check
		results, err = queries.FilterByTeam(ctx, teamF)
	} else {
		// table scan - slow!
		results, err = queries.All(ctx)
	}
	return
}

// resultsOut converts the structs to map for output
func resultsOut(results []ghs.GithubStandard, response *resp.Response) (rows []map[string]interface{}, base int, ext int) {
	base = 0
	ext = 0
	rows = []map[string]interface{}{}
	for _, item := range results {
		base += item.CompliantBaseline
		ext += item.CompliantExtended

		if m, err := convert.Map(item); err == nil {
			rows = append(rows, m)
		} else {
			slog.Error("error converting result to map", slog.String("err", err.Error()))
		}
	}
	slog.Info("result out", slog.Int("r", len(response.Errors)))
	return
}

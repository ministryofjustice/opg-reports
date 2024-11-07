package releasesapi

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/releases"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesdb"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesio"
)

const (
	Segment string = "releases"
	Tag     string = "Releases"
)

var teamsListingDescription string = `Returns list of all teams within the database.`

// apiTeamsListingHandler lists all teams in the database
func apiTeamsListingHandler(ctx context.Context, input *releasesio.ReleasesTeamsInput) (response *releasesio.ReleaseTeamsOutput, err error) {
	var (
		result     []*releases.Team
		db         *sqlx.DB
		dbFilepath string                         = ctx.Value(Segment).(string)
		stmt       datastore.NamedSelectStatement = releasesdb.ListTeams
		bdy        *releasesio.ReleasesTeamsBody  = &releasesio.ReleasesTeamsBody{
			Request: input,
			Type:    "teams-list",
		}
	)
	// grab db connection
	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	if err != nil {
		return
	}
	defer db.Close()

	// select them all
	if result, err = datastore.SelectMany[*releases.Team](ctx, db, stmt, input); err == nil {
		bdy.Result = result
	}

	response = &releasesio.ReleaseTeamsOutput{
		Body: bdy,
	}

	return
}

// // apiReleasesListingHandler lists all releases without any grouping all filters
// func apiReleasesListingHandler(ctx context.Context, input *releasesio.ReleasesListAllInput) (response *releasesio.ReleaseListAllOutput, err error) {
// 	var (
// 		result     []*releases.Release
// 		db         *sqlx.DB
// 		dbFilepath string                          = ctx.Value(Segment).(string)
// 		stmt       datastore.NamedSelectStatement  = releasesdb.ListReleases
// 		bdy        *releasesio.ReleasesListAllBody = &releasesio.ReleasesListAllBody{
// 			Request: input,
// 			Type:    "releases-list",
// 		}
// 	)
// 	// grab db connection
// 	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
// 	if err != nil {
// 		return
// 	}
// 	defer db.Close()

// 	// select them all
// 	if result, err = datastore.SelectMany[*releases.Release](ctx, db, stmt, input); err == nil {
// 		bdy.Result = result
// 	}

// 	response = &releasesio.ReleaseListAllOutput{
// 		Body: bdy,
// 	}

// 	return
// }

var releasesIntervalDescription = `Returns count of releases per interval between the start and end date passed.`

// apiReleasesIntervalHandler lists all releases between start and end date, grouped only by time period
func apiReleasesIntervalHandler(ctx context.Context, input *releasesio.ReleasesInput) (response *releasesio.ReleaseOutput, err error) {
	var (
		result     []*releases.Release
		db         *sqlx.DB
		dbFilepath string                         = ctx.Value(Segment).(string)
		stmt       datastore.NamedSelectStatement = releasesdb.ListReleasesGroupedByInterval
		bdy        *releasesio.ReleasesBody       = &releasesio.ReleasesBody{
			Request:   input,
			Type:      "releases",
			DateRange: convert.DateRange(input.StartTime(), input.EndTime(), input.Interval),
		}
	)
	// grab db connection
	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	if err != nil {
		return
	}
	defer db.Close()

	if input.Unit != "" {
		stmt = releasesdb.ListReleasesGroupedByIntervalFilter
	}

	// select them all
	if result, err = datastore.SelectMany[*releases.Release](ctx, db, stmt, input); err == nil {
		bdy.Result = result
	}

	response = &releasesio.ReleaseOutput{
		Body: bdy,
	}

	return
}

var releasesIntervalTeamDescription = `Returns count of releases per interval between the start and end date passed grouped by time period and team names.

As more than one team can be attached to a release the total values per time period may look higher.
`

// apiReleasesIntervalTeamHandler lists all releases between start and end date, grouped by time period and teams
func apiReleasesIntervalTeamHandler(ctx context.Context, input *releasesio.ReleasesInput) (response *releasesio.ReleaseOutput, err error) {
	var (
		result     []*releases.Release
		db         *sqlx.DB
		dbFilepath string                         = ctx.Value(Segment).(string)
		stmt       datastore.NamedSelectStatement = releasesdb.ListReleasesGroupedByIntervalAndTeam
		bdy        *releasesio.ReleasesBody       = &releasesio.ReleasesBody{
			Request:     input,
			Type:        "releases-unit",
			ColumnOrder: []string{"unit"},
			DateRange:   convert.DateRange(input.StartTime(), input.EndTime(), input.Interval),
		}
	)
	// grab db connection
	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	if err != nil {
		return
	}
	defer db.Close()

	if input.Unit != "" {
		stmt = releasesdb.ListReleasesGroupedByIntervalAndTeamFilter
	}

	// select them all
	if result, err = datastore.SelectMany[*releases.Release](ctx, db, stmt, input); err == nil {
		bdy.Result = result
		bdy.ColumnValues = datastore.ColumnValues(result, bdy.ColumnOrder)
	}

	response = &releasesio.ReleaseOutput{
		Body: bdy,
	}

	return
}

func Register(api huma.API) {
	var uri string

	uri = "/{version}/releases/github/teams"
	slog.Info("[releasseapi.Register] ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   "get-releases-teams-list",
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List all teams",
		Description:   teamsListingDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiTeamsListingHandler)

	// uri = "/{version}/releases/github/all"
	// slog.Info("[releasseapi.Register] ", slog.String("uri", uri))
	// huma.Register(api, huma.Operation{
	// 	OperationID:   "get-releases-list",
	// 	Method:        http.MethodGet,
	// 	Path:          uri,
	// 	Summary:       "List all releases",
	// 	Description:   teamsListingDescription,
	// 	DefaultStatus: http.StatusOK,
	// 	Tags:          []string{Tag},
	// }, apiReleasesListingHandler)

	uri = "/{version}/releases/github/{start_date}/{end_date}/{interval}"
	slog.Info("[releasseapi.Register] ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   "get-releases",
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Releases",
		Description:   releasesIntervalDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiReleasesIntervalHandler)

	uri = "/{version}/releases/github/unit/{start_date}/{end_date}/{interval}"
	slog.Info("[releasseapi.Register] ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   "get-releases-unit",
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Releases per unit",
		Description:   releasesIntervalTeamDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiReleasesIntervalTeamHandler)

}

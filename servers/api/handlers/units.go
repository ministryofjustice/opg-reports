package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

var (
	UnitsSegment string   = "unit"
	UnitTags     []string = []string{"Units"}
)

const UnitsListDescription string = `Returns all units within the database.`
const UnitListOperationID string = "get-units-list"
const unitListSQL string = `
SELECT
	units.*,
	json_group_array(
		DISTINCT json_object(
			'id', aws_accounts.id,
			'number', aws_accounts.number,
			'name', aws_accounts.name,
			'label', aws_accounts.label,
			'environment', aws_accounts.environment
		)
	) filter ( where aws_accounts.id is not null) as aws_accounts,
	json_group_array(
		DISTINCT json_object(
			'id', github_teams.id,
			'slug', github_teams.slug
		)
	) filter ( where github_teams.id is not null) as github_teams
FROM units
LEFT JOIN aws_accounts ON aws_accounts.unit_id = units.id
LEFT JOIN github_teams_units on github_teams_units.unit_id = units.id
LEFT JOIN github_teams on github_teams.id = github_teams_units.github_team_id
GROUP BY units.id
ORDER BY units.name ASC;
`

// ApiUnitsListHandler queries the database for all units and returns them as a list including
// joins with github teams and aws accounts. There is no option to filter of limit the results.
//
// Endpoints:
//
//	/version/units/list
func ApiUnitsListHandler(ctx context.Context, input *inout.VersionInput) (response *inout.UnitsListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.Unit       = []*models.Unit{}
		dbPath  string               = ctx.Value(dbPathKey).(string)
		body    *inout.UnitsListBody = inout.NewUnitsListBody()
	)
	body.Request = input
	body.Operation = UnitListOperationID
	// setup response
	response = &inout.UnitsListResponse{}
	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] units list adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.Unit](ctx, adaptor, unitListSQL, nil)
	if err != nil {
		slog.Error("[api] units list select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("units list selection failed."))
	} else {
		body.Result = results
	}
	response.Body = body
	return
}

// Register attaches the handler to the main api
func RegisterUnits(api huma.API) {
	var uri string = "/{version}/" + UnitsSegment + "/list"

	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   UnitListOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List units",
		Description:   UnitsListDescription,
		DefaultStatus: http.StatusOK,
		Tags:          UnitTags,
	}, ApiUnitsListHandler)

}

// Package `api` handles returning account information.
//
// Will return all accounts or those belonging to the team
// name passed into via the http request
package api

import (
	"opg-reports/report/internal/domains/account/types"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/handler"
	"opg-reports/report/packages/types/models"
)

const label string = `get-accounts`

const stmt string = `
SELECT
	id,
	name,
	label,
	environment,
	team_name as team
FROM accounts
WHERE
	{TEAM_FILTER}
	id IS NOT NULL
ORDER BY
	team_name,
	name,
	environment ASC
;`

func Config(opts *args.API) *handler.ApiConfig {
	return &handler.ApiConfig{
		Name:     label,
		Database: opts.DB,
		Selector: &types.Select{Statement: stmt},
	}
}

func Response(opts *args.API) *models.ApiResponse {
	return &models.ApiResponse{
		Versions: opts.Versions,
	}
}

func T() *types.Account {
	return &types.Account{}
}

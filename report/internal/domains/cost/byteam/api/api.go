// Package `api` returns cost data grouped by team name
package api

const label string = `get-costs-all`

const stmt string = `
SELECT
	costs.month as month,
	CAST(COALESCE(SUM(cost), 0) as REAL) as cost,
	IIF(accounts.team_name != "", accounts.team_name, "")  as team
FROM costs
LEFT JOIN accounts on accounts.id = costs.account_id
WHERE
	{TEAM_FILTER}
	costs.service != 'Tax'
	AND costs.month IN (:months)
GROUP BY
	costs.month,
	accounts.team_name
ORDER BY
	accounts.team_name ASC
;`

// // SETUP vars
// var (
// 	request    *models.Request
// 	statement  *types.Select // the sql statement
// 	results    []*types.CostByTeam
// 	response   *models.APIResponse[*types.CostByTeam, *models.Request, *models.Filter] // the response type
// 	apiHandler *handler.APIConfig[*types.CostByTeam, *models.Request, *models.Filter]  // the api handler
// )

// // tabulateF is used to update the results and convert items to a table
// func tabulateF(resp *models.APIResponse[*types.CostByTeam, *models.Request, *models.Filter]) {

// }

// // Handler returns the api handler config to be used to fetch data.
// func Handler(opts *args.API) *handler.APIConfig[*types.CostByTeam, *models.Request, *models.Filter] {
// 	request = &models.Request{}
// 	statement = &types.Select{Statement: stmt}
// 	results = []*types.CostByTeam{}

// 	// the parent response
// 	response = &models.APIResponse[*types.CostByTeam, *models.Request, *models.Filter]{
// 		Version: opts.Versions.Version,
// 		SHA:     opts.Versions.SHA,
// 	}

// 	apiHandler = &handler.APIConfig[*types.CostByTeam, *models.Request, *models.Filter]{
// 		Label:     label,
// 		DB:        opts.DB,
// 		Statement: statement,
// 		Request:   request,
// 		Results:   results,
// 		Response:  response,
// 	}

// 	return apiHandler
// }

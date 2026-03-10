package api

const label string = `get-teams`

const stmt string = `
SELECT
	name
FROM teams
WHERE
	name != 'legacy'
	AND name != 'org'
ORDER BY
	name asc
;`

// // SETUP vars
// var (
// 	request    *models.Request
// 	statement  *types.Select
// 	results    []*types.Team
// 	response   *models.APIResponse[*types.Team, *types.Team, *models.Request, *models.Filter] // the response type
// 	apiHandler *handler.APIConfig[*types.Team, *types.Team, *models.Request, *models.Filter]  // the api handler
// )

// // Handler returns the api handler config to be used to fetch data.
// func Handler(opts *args.API) *handler.APIConfig[*types.Team, *types.Team, *models.Request, *models.Filter] {
// 	request = &models.Request{}
// 	statement = &types.Select{Statement: stmt}
// 	results = []*types.Team{}
// 	response = &models.APIResponse[*types.Team, *types.Team, *models.Request, *models.Filter]{
// 		Version: opts.Versions.Version,
// 		SHA:     opts.Versions.SHA,
// 	}

// 	apiHandler = &handler.APIConfig[*types.Team, *types.Team, *models.Request, *models.Filter]{
// 		Label:     label,
// 		DB:        opts.DB,
// 		Statement: statement,
// 		Request:   request,
// 		Results:   results,
// 		Response:  response,
// 	}

// 	return apiHandler
// }

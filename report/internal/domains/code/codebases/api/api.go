package api

const label string = `get-code-all`

const stmt string = `
SELECT
	full_name,
	name,
	url,
	archived
FROM codebases
WHERE
	archived = 0
ORDER BY
	name ASC
;`

// // SETUP vars
// var (
// 	request    *models.Request
// 	statement  *types.Select // the sql statement
// 	results    []*types.Codebase
// 	response   *models.APIResponse[*types.Codebase, *types.Codebase, *models.Request, *models.Filter] // the response type
// 	apiHandler *handler.APIConfig[*types.Codebase, *types.Codebase, *models.Request, *models.Filter]  // the api handler
// )

// // Handler returns the api handler config to be used to fetch data.
// func Handler(opts *args.API) *handler.APIConfig[*types.Codebase, *types.Codebase, *models.Request, *models.Filter] {
// 	request = &models.Request{}
// 	statement = &types.Select{Statement: stmt}
// 	results = []*types.Codebase{}

// 	response = &models.APIResponse[*types.Codebase, *types.Codebase, *models.Request, *models.Filter]{
// 		Version: opts.Versions.Version,
// 		SHA:     opts.Versions.SHA,
// 	}

// 	apiHandler = &handler.APIConfig[*types.Codebase, *types.Codebase, *models.Request, *models.Filter]{
// 		Label:     label,
// 		DB:        opts.DB,
// 		Statement: statement,
// 		Request:   request,
// 		Results:   results,
// 		Response:  response,
// 	}

// 	return apiHandler
// }

package api

const stmtGithubCodeOwnerInsert string = `
INSERT INTO github_codeownership (
	codeowner,
	repository,
	team
) VALUES (
	:codeowner,
	:repository,
	:team
) ON CONFLICT (codeowner,repository,team)
 	DO UPDATE SET team=excluded.team
RETURNING id;`

// stmtGithubCodeOwnerSelectAll is sql used to fetch all teams and the join to aws accounts
const stmtGithubCodeOwnerSelectAll string = `
SELECT
	codeowner,
	repository,
	json_object(
		'name', github_codeownership.team
	) as team
FROM github_codeownership
LEFT JOIN teams ON github_codeownership.team = teams.name
GROUP BY codeowner
ORDER BY codeowner ASC;`

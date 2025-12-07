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
	team
FROM github_codeownership
ORDER BY codeowner ASC;`

const stmtGithubCodeOwnerSelectForTeam string = `
SELECT
	codeowner,
	repository,
	team
FROM github_codeownership
WHERE
	lower(team)=lower(:team)
ORDER BY codeowner ASC;`

const stmtGithubCodeOwnerTruncate string = `DELETE FROM github_codeownership;`

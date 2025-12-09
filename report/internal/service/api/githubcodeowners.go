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

// stmtGithubCodeOwnerSelectAll is sql used to fetch all codeowners
const stmtGithubCodeOwnerSelectAll string = `
SELECT
	codeowner,
	repository,
	team
FROM github_codeownership
ORDER BY codeowner ASC;`

// stmtGithubCodeOwnerSelectForTeam is sql used to filter results based on team name
const stmtGithubCodeOwnerSelectForTeam string = `
SELECT
	codeowner,
	repository,
	team
FROM github_codeownership
WHERE
	lower(team)=lower(:team)
ORDER BY codeowner ASC;`

// stmtGithubCodeOwnerSelectForTeam is sql used to filter results based on codeowner
const stmtGithubCodeOwnerSelectForCodeOwner string = `
SELECT
	codeowner,
	repository,
	team
FROM github_codeownership
WHERE
	lower(codeowner)=lower(:codeowner)
ORDER BY codeowner ASC;`

const stmtGithubCodeOwnerTruncate string = `DELETE FROM github_codeownership;`

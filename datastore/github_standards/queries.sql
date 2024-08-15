-- name: Insert :one
INSERT INTO github_standards(
    ts,
    default_branch,
    full_name,
    name,
    owner,
    license,
    last_commit_date,
    created_at,
    count_of_clones,
    count_of_forks,
    count_of_pull_requests,
    count_of_web_hooks,
    has_code_of_conduct,
    has_codeowner_approval_required,
    has_contributing_guide,
    has_default_branch_of_main,
    has_default_branch_protection,
    has_delete_branch_on_merge,
    has_description,
    has_discussions,
    has_downloads,
    has_issues,
    has_license,
    has_pages,
    has_pull_request_approval_required,
    has_readme,
    has_rules_enforced_for_admins,
    has_vulnerability_alerts,
    has_wiki,
    is_archived,
    is_private,
    compliant_baseline,
    compliant_extended,
    teams
) VALUES (
    ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?
) RETURNING id;

-- name: All :many
SELECT * FROM github_standards
ORDER BY name, created_at ASC;

-- name: Count :one
SELECT count(*) FROM github_standards;

-- name: TotalCountCompliantBaseline :one
SELECT count(*) FROM github_standards
WHERE compliant_baseline=1;

-- name: TotalCountCompliantExtended :one
SELECT count(*) FROM github_standards
WHERE compliant_extended=1;

-- name: FilterByIsArchived :many
SELECT * FROM github_standards
WHERE is_archived = ?
ORDER BY name, created_at ASC;

-- name: FilterByTeam :many
SELECT * FROM github_standards
WHERE teams LIKE ?
ORDER BY name, created_at ASC;

-- name: FilterByIsArchivedAndTeam :many
SELECT * FROM github_standards
WHERE
    is_archived = ? AND
    teams LIKE ?
ORDER BY name, created_at ASC;

-- name: Track :exec
INSERT INTO github_standards_tracker (run_date) VALUES(?) ;

-- name: Age :one
SELECT run_date FROM github_standards_tracker LIMIT 1;

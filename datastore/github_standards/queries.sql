-- name: All :many
SELECT * FROM github_standards
ORDER BY name, created_at ASC;

-- name: ArchivedFilter :many
SELECT * FROM github_standards
WHERE is_archived = ?
ORDER BY name, created_at ASC;

-- name: TeamFilter :many
SELECT * FROM github_standards
WHERE teams LIKE ?
ORDER BY name, created_at ASC;

-- name: ArchivedTeamFilter :many
SELECT * FROM github_standards
WHERE
    is_archived = ? AND
    teams LIKE ?
ORDER BY name, created_at ASC;

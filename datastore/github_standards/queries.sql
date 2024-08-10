-- name: All :many
SELECT * FROM github_standards
ORDER BY name, created_at ASC;

-- name: Archived :many
SELECT * FROM github_standards
WHERE is_archived = ?
ORDER BY name, created_at ASC;

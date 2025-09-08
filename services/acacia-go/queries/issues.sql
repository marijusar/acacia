-- name: GetAllIssues :many
SELECT * FROM issues
ORDER BY created_at DESC;

-- name: GetIssueByID :one
SELECT * FROM issues
WHERE id = $1;

-- name: CreateIssue :one
INSERT INTO issues (name, description, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;

-- name: UpdateIssue :one
UPDATE issues
SET name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteIssue :one
DELETE FROM issues
WHERE id = $1
RETURNING *;
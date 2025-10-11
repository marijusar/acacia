-- name: GetIssuesByColumnId :many
SELECT
    *
FROM
    issues
WHERE
    column_id = $1
ORDER BY
    created_at DESC;

-- name: GetIssueByID :one
SELECT
    *
FROM
    issues
WHERE
    id = $1;

-- name: CreateIssue :one
INSERT INTO issues (name, column_id, description, created_at, updated_at)
    VALUES ($1, $2, $3, NOW(), NOW())
RETURNING
    *;

-- name: UpdateIssue :one
UPDATE
    issues
SET
    name = COALESCE(@name, name),
    description = COALESCE(@description, description),
    column_id = COALESCE(@column_id, column_id),
    updated_at = NOW()
WHERE
    id = @id
RETURNING
    *;

-- name: ReassignAllIssuesFromColumn :exec
UPDATE
    issues
SET
    column_id = @target_column
WHERE
    column_id = @source_column;

-- name: DeleteIssue :exec
DELETE FROM issues
WHERE id = $1;


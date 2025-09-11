-- name: GetAllProjectStatusColumns :many
SELECT * FROM project_status_columns
ORDER BY project_id, position_index;

-- name: GetProjectStatusColumnByID :one
SELECT * FROM project_status_columns
WHERE id = $1;

-- name: GetProjectStatusColumnsByProjectID :many
SELECT * FROM project_status_columns
WHERE project_id = $1
ORDER BY position_index;

-- name: CreateProjectStatusColumn :one
INSERT INTO project_status_columns (project_id, name, position_index, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING *;

-- name: UpdateProjectStatusColumn :one
UPDATE project_status_columns
SET name = $2,
    position_index = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProjectStatusColumn :one
DELETE FROM project_status_columns
WHERE id = $1
RETURNING *;
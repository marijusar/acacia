-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = $1;

-- name: CreateProject :one
INSERT INTO projects (name, created_at, updated_at)
VALUES ($1, NOW(), NOW())
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET name = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProject :one
DELETE FROM projects
WHERE id = $1
RETURNING *;

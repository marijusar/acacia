-- name: GetProjects :many
SELECT
    p.*
FROM
    projects p
    JOIN team_members tm ON p.team_id = tm.team_id
WHERE
    tm.user_id = $1;

-- name: GetProjectByID :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

-- name: CreateProject :one
INSERT INTO projects (name, team_id, created_at, updated_at)
    VALUES ($1, $2, NOW(), NOW())
RETURNING
    *;

-- name: UpdateProject :one
UPDATE
    projects
SET
    name = $2,
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteProject :one
DELETE FROM projects
WHERE id = $1
RETURNING
    *;

-- name: GetProjectIssues :many
SELECT
    issues.*
FROM
    project_status_columns
    JOIN issues ON project_status_columns.id = issues.column_id
WHERE
    project_id = $1;


-- name: CreateTeam :one
INSERT INTO teams (name, description, created_at, updated_at)
    VALUES ($1, $2, NOW(), NOW())
RETURNING *;

-- name: GetTeamByID :one
SELECT * FROM teams WHERE id = $1;

-- name: AddTeamMember :one
INSERT INTO team_members (team_id, user_id, joined_at)
    VALUES ($1, $2, NOW())
RETURNING *;

-- name: RemoveTeamMember :exec
DELETE FROM team_members WHERE team_id = $1 AND user_id = $2;

-- name: GetTeamMembers :many
SELECT u.* FROM users u
JOIN team_members tm ON u.id = tm.user_id
WHERE tm.team_id = $1;

-- name: GetUserTeams :many
SELECT t.* FROM teams t
JOIN team_members tm ON t.id = tm.team_id
WHERE tm.user_id = $1;

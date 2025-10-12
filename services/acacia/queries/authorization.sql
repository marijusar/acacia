-- name: CheckUserTeamMembership :one
SELECT EXISTS (
    SELECT 1
    FROM team_members
    WHERE team_id = $1 AND user_id = $2
) AS is_member;

-- name: GetTeamIDByProject :one
SELECT team_id
FROM projects
WHERE id = $1;

-- name: GetTeamIDByProjectStatusColumn :one
SELECT p.team_id
FROM project_status_columns psc
JOIN projects p ON psc.project_id = p.id
WHERE psc.id = $1;

-- name: GetTeamIDByIssue :one
SELECT p.team_id
FROM issues i
JOIN project_status_columns psc ON i.column_id = psc.id
JOIN projects p ON psc.project_id = p.id
WHERE i.id = $1;

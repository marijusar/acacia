-- name: CreateTeamLLMAPIKey :one
INSERT INTO teams_llm_api_keys (team_id, provider, encrypted_key)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTeamLLMAPIKeyByTeamID :one
SELECT * FROM teams_llm_api_keys
WHERE team_id = $1 AND provider = $2 AND is_active = true
LIMIT 1;

-- name: GetTeamLLMAPIKeyByProjectID :one
SELECT tlak.* FROM teams_llm_api_keys tlak
JOIN projects p ON p.team_id = tlak.team_id
WHERE p.id = $1 AND tlak.provider = $2 AND tlak.is_active = true
LIMIT 1;

-- name: UpdateTeamLLMAPIKey :one
UPDATE teams_llm_api_keys
SET encrypted_key = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTeamLLMAPIKey :exec
DELETE FROM teams_llm_api_keys
WHERE id = $1;

-- name: CheckTeamLLMAPIKeyExists :one
SELECT EXISTS(
    SELECT 1 FROM teams_llm_api_keys tlak
    JOIN projects p ON p.team_id = tlak.team_id
    WHERE p.id = $1 AND tlak.provider = $2 AND tlak.is_active = true
) AS exists;

-- name: UpdateLastUsedAt :exec
UPDATE teams_llm_api_keys
SET last_used_at = NOW()
WHERE id = $1;

-- name: GetAllTeamLLMAPIKeys :many
SELECT * FROM teams_llm_api_keys
WHERE team_id = $1 AND is_active = true
ORDER BY created_at DESC;

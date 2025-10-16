-- name: CreateConversation :one
INSERT INTO conversations (user_id, title, provider, model)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: GetConversationByID :one
SELECT
    *
FROM
    conversations
WHERE
    id = $1;

-- name: GetConversationsByUser :many
SELECT
    *
FROM
    conversations
WHERE
    user_id = $1
ORDER BY
    created_at DESC;

-- name: GetLatestConversationByUser :one
SELECT
    *
FROM
    conversations
WHERE
    user_id = $1
ORDER BY
    created_at DESC
LIMIT 1;

-- name: UpdateConversationTitle :one
UPDATE
    conversations
SET
    title = $2,
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteConversation :exec
DELETE FROM conversations
WHERE id = $1;


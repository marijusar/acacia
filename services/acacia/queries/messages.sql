-- name: CreateMessage :one
INSERT INTO messages (
    conversation_id,
    role,
    content,
    sequence_number
) VALUES (
    $1, $2, $3,
    (SELECT COALESCE(MAX(sequence_number), 0) + 1 FROM messages WHERE conversation_id = $1)
) RETURNING *;

-- name: GetMessagesByConversationID :many
SELECT * FROM messages
WHERE conversation_id = $1
ORDER BY sequence_number ASC;

-- name: GetMessageByID :one
SELECT * FROM messages
WHERE id = $1;

-- name: DeleteMessageByID :exec
DELETE FROM messages
WHERE id = $1;

-- name: DeleteMessagesByConversationID :exec
DELETE FROM messages
WHERE conversation_id = $1;

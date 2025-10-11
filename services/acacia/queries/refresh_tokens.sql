-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, jti, expires_at, created_at)
    VALUES ($1, $2, $3, NOW())
RETURNING
    *;

-- name: GetRefreshTokenByUserAndJTI :one
SELECT
    *
FROM
    refresh_tokens
WHERE
    user_id = $1
    AND jti = $2
    AND expires_at > NOW()
    AND revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE
    refresh_tokens
SET
    revoked_at = NOW()
WHERE
    user_id = $1
    AND jti = $2;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE
    refresh_tokens
SET
    revoked_at = NOW()
WHERE
    user_id = $1
    AND revoked_at IS NULL;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();

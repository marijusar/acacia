CREATE TABLE IF NOT EXISTS refresh_tokens (
    id bigserial PRIMARY KEY,
    user_id bigint REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    jti varchar(255) NOT NULL UNIQUE,
    expires_at timestamp NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    revoked_at timestamp
);

CREATE INDEX idx_refresh_tokens_user_id_jti ON refresh_tokens(user_id, jti);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

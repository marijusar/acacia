CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email varchar(255) NOT NULL UNIQUE,
    name varchar(100) NOT NULL,
    password_hash varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

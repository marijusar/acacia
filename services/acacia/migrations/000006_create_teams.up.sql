CREATE TABLE IF NOT EXISTS teams (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    description text,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_teams_created_at ON teams(created_at);

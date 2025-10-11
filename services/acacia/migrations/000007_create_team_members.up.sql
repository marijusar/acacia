CREATE TABLE IF NOT EXISTS team_members (
    id bigserial PRIMARY KEY,
    team_id bigint NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at timestamp NOT NULL DEFAULT NOW(),
    UNIQUE(team_id, user_id)
);

CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);
CREATE INDEX idx_team_members_team_user ON team_members(team_id, user_id);

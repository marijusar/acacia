CREATE TABLE IF NOT EXISTS teams_llm_api_keys (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    encrypted_key TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMP,
    UNIQUE(team_id, provider)
);

CREATE INDEX idx_teams_llm_api_keys_team_id ON teams_llm_api_keys(team_id);
CREATE INDEX idx_teams_llm_api_keys_provider ON teams_llm_api_keys(provider);
CREATE INDEX idx_teams_llm_api_keys_team_provider ON teams_llm_api_keys(team_id, provider);

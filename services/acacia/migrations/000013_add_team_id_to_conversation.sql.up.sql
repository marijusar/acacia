ALTER TABLE conversations
    ADD COLUMN team_id bigint NOT NULL REFERENCES teams (id);


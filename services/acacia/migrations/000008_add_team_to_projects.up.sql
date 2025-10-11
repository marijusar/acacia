ALTER TABLE projects
ADD COLUMN team_id bigint REFERENCES teams(id) ON DELETE CASCADE;

CREATE INDEX idx_projects_team_id ON projects(team_id);

-- Note: After this migration, existing projects will have NULL team_id
-- You may want to create a default team and assign existing projects to it
-- Or make team_id NOT NULL after handling existing data

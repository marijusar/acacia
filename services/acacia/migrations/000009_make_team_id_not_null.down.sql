-- Revert team_id to allow NULL values

ALTER TABLE projects
ALTER COLUMN team_id DROP NOT NULL;

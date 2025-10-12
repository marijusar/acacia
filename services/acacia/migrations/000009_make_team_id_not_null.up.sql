-- Make team_id NOT NULL in projects table
-- This ensures every project must belong to a team

-- First, delete any projects without a team_id (orphaned projects)
DELETE FROM projects WHERE team_id IS NULL;

-- Now make the column NOT NULL
ALTER TABLE projects
ALTER COLUMN team_id SET NOT NULL;

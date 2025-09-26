-- Add the column_id foreign key column to issues table
ALTER TABLE issues
    ADD COLUMN column_id bigint NOT NULL REFERENCES project_status_columns (id);

-- Create index on position_index for better query performance
CREATE INDEX idx_project_status_columns_position_index ON project_status_columns (position_index);

-- Optional: Create composite index for project_id + position_index for even better performance
-- when querying columns within a project ordered by position
CREATE INDEX idx_project_status_columns_project_position ON project_status_columns (project_id, position_index);

-- Optional: Create index on column_id in issues table for better join performance
CREATE INDEX idx_issues_column_id ON issues (column_id);


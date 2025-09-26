DROP INDEX IF EXISTS idx_issues_column_id;

DROP INDEX IF EXISTS idx_project_status_columns_project_position;

DROP INDEX IF EXISTS idx_project_status_columns_position_index;

-- Drop the foreign key column
ALTER TABLE issues
    DROP COLUMN IF EXISTS column_id;


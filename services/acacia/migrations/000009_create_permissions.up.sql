-- Permissions table: defines all possible permissions in the system
CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    name varchar(100) NOT NULL UNIQUE, -- e.g., 'project:view', 'project:edit'
    resource_type varchar(50) NOT NULL, -- e.g., 'project', 'team', 'issue'
    action varchar(50) NOT NULL, -- e.g., 'view', 'edit', 'delete'
    description text,
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_permissions_resource_action ON permissions(resource_type, action);

-- Insert default permissions
INSERT INTO permissions (name, resource_type, action, description) VALUES
-- Project permissions
('project:view', 'project', 'view', 'View project details'),
('project:edit', 'project', 'edit', 'Edit project details'),
('project:delete', 'project', 'delete', 'Delete project'),
('project:create', 'project', 'create', 'Create new project'),

-- Team permissions
('team:view', 'team', 'view', 'View team details'),
('team:edit', 'team', 'edit', 'Edit team details'),
('team:delete', 'team', 'delete', 'Delete team'),
('team:manage_members', 'team', 'manage_members', 'Add/remove team members'),

-- Issue permissions
('issue:view', 'issue', 'view', 'View issues'),
('issue:create', 'issue', 'create', 'Create issues'),
('issue:edit', 'issue', 'edit', 'Edit issues'),
('issue:delete', 'issue', 'delete', 'Delete issues'),

-- Column permissions
('column:view', 'column', 'view', 'View columns'),
('column:create', 'column', 'create', 'Create columns'),
('column:edit', 'column', 'edit', 'Edit columns'),
('column:delete', 'column', 'delete', 'Delete columns');

CREATE TABLE IF NOT EXISTS projects (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_status_columns (
    id bigserial PRIMARY KEY,
    project_id integer REFERENCES projects (id) NOT NULL,
    name varchar(255) NOT NULL,
    position_index smallint NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);


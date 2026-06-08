CREATE TYPE project_permission AS ENUM (
    'read',
    'edit'
);

CREATE TABLE project_collaborators (
    id UUID PRIMARY KEY,

    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,

    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    permission project_permission NOT NULL DEFAULT 'read',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(project_id, user_id)
);
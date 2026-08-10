CREATE TABLE projects (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT projects_owner_name_unique UNIQUE (owner_id, name),
    CONSTRAINT projects_owner_slug_unique UNIQUE (owner_id, slug)
);

CREATE INDEX projects_owner_id_idx
    ON projects (owner_id);

CREATE INDEX projects_created_at_idx
    ON projects (created_at DESC);
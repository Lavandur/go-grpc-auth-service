CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY,
    login VARCHAR NOT NULL,
    visible_id VARCHAR NOT NULL,
    hashed_password VARCHAR NOT NULL,
    person JSONB NOT NULL,
    role_ids VARCHAR[],
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP,
    last_password_restore_date TIMESTAMP,
    search_index VARCHAR
);

CREATE UNIQUE INDEX unique_login ON users (login) WHERE deleted_date IS NULL;
CREATE UNIQUE INDEX unique_visible_id ON users (visible_id) WHERE deleted_date IS NULL;
CREATE INDEX ON users (search_index);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id UUID NOT NULL,
  permission VARCHAR NOT NULL,

  CONSTRAINT unique_role_permission UNIQUE (role_id, permission)
);

CREATE TABLE IF NOT EXISTS permissions (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description JSONB
);

CREATE TABLE IF NOT EXISTS roles (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description JSONB,
    created_at TIMESTAMP
);

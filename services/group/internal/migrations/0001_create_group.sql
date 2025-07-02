-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE SCHEMA IF NOT EXISTS "group";
SET search_path TO "group";

CREATE TABLE IF NOT EXISTS groups
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT UNIQUE NOT NULL,
    department TEXT,
    year       INTEGER,
    created_at timestamptz      DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_groups
(
    user_id  UUID NOT NULL,
    group_id UUID NOT NULL,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES "user".users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

-- +goose Down
SET search_path TO "group";

DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS groups;
DROP SCHEMA IF EXISTS "group" CASCADE;
DROP EXTENSION IF EXISTS "pgcrypto";

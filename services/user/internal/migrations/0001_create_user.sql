-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE SCHEMA IF NOT EXISTS "user";
SET search_path TO "user";

DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_status') THEN CREATE TYPE role_status AS ENUM ('student', 'teacher', 'admin'); END IF; END $$;

CREATE TABLE IF NOT EXISTS roles
(
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name     TEXT        NOT NULL,
    barcode  TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL,
    role_id  INT REFERENCES roles (id)
);

CREATE TABLE IF NOT EXISTS permissions
(
    id     SERIAL PRIMARY KEY,
    action TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions
(
    role_id       INT REFERENCES roles (id) ON DELETE CASCADE,
    permission_id INT REFERENCES permissions (id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

INSERT INTO roles (name)
VALUES ('admin'),
       ('teacher'),
       ('student')
ON CONFLICT DO NOTHING;

-- +goose Down
SET search_path TO "user";

DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TYPE IF EXISTS role_status;
DROP SCHEMA IF EXISTS "user" CASCADE;
DROP EXTENSION IF EXISTS "pgcrypto";

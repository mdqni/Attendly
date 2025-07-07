CREATE SCHEMA IF NOT EXISTS "auth";

SET search_path TO "auth";

CREATE TABLE IF NOT EXISTS roles
(
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY,
    name     TEXT        NOT NULL,
    barcode  TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL,
    role_id  INT         REFERENCES roles (id) ON DELETE SET NULL
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

CREATE TABLE IF NOT EXISTS refresh_tokens
(
    token      TEXT PRIMARY KEY,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);

INSERT INTO roles (name)
VALUES ('admin'),
       ('teacher'),
       ('student')
ON CONFLICT DO NOTHING;

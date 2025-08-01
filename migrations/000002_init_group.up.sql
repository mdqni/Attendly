CREATE SCHEMA IF NOT EXISTS "group";
SET search_path TO "group";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";


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
    FOREIGN KEY (user_id) REFERENCES "auth".users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

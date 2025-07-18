CREATE SCHEMA IF NOT EXISTS "user";
SET search_path TO "user";

CREATE TABLE IF NOT EXISTS user_profiles
(
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT,
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (id) REFERENCES "auth".users (id) ON DELETE CASCADE
);

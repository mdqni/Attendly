CREATE SCHEMA IF NOT EXISTS "userschema";
SET search_path TO "userschema";

CREATE TABLE IF NOT EXISTS user_profiles
(
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT,
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (id) REFERENCES "auth".users (id) ON DELETE CASCADE
);

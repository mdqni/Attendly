SET search_path TO "userschema";

DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TYPE IF EXISTS role_status;
DROP SCHEMA IF EXISTS "user" CASCADE;
DROP EXTENSION IF EXISTS "pgcrypto";

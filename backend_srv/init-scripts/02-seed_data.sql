-- NOTE dont create database here ==> database has been created at docker initiation
-- because it can CRASH the docker

-- extension
CREATE EXTENSION IF NOT EXISTS postgis;


-- 1. Create users Table First (The Parent)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    "password" VARCHAR(100) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);
-- index username
CREATE INDEX username_idx ON users (username);

INSERT INTO users (
    username,
    "password"
) VALUES (
    'admin01',
    'admin123'
);




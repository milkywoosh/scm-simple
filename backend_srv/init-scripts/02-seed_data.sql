-- NOTE dont create database here ==> database has been created at docker initiation
-- because it can CRASH the docker

-- extension
CREATE EXTENSION IF NOT EXISTS postgis;


-- 1. Create users Table First (The Parent)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    "password" VARCHAR(100) UNIQUE,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    updated_at DATE
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

-- constraint: 1 warehouse only have 1 parent, 1 parent have multi child
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    location_code VARCHAR(60) NOT NULL UNIQUE,
    "description" VARCHAR(100),
    parent_location_code VARCHAR(60),
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    updated_at DATE,
    location_point geography(Point, 4326)
);
-- enforce SRID 4326
ALTER TABLE warehouses ADD CONSTRAINT enforce_srid CHECK (ST_SRID(location_point) = 4326);
-- index location code
CREATE INDEX location_code_idx ON warehouses (location_code);

INSERT INTO warehouses (
    location_code,
    "description",
    parent_location_code
) VALUES (
    'PGC001',
    'pusat grosir cililitan',
    null
);


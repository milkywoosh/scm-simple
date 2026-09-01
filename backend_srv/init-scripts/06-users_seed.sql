
-- 

CREATE TABLE users (
    id serial primary key,
    username varchar(60) unique, 
    hashed_password varchar(100),
    full_name varchar(100),
    email varchar(90) UNIQUE,
    password_changed_at TIMESTAMPTZ default null,
    created_at TIMESTAMPTZ
);

CREATE INDEX idx_users_username_email
    ON users (username, email);

INSERT INTO users (username, hashed_password, full_name, email, password_changed_at, created_at) 
VALUES('ben01', NULL, 'ben one', 'ben@gmail.com', now(), now());

INSERT INTO users (username, hashed_password, full_name, email, password_changed_at, created_at) 
VALUES('ron01', NULL, 'ron one', 'ron@gmail.com', now(), now());

INSERT INTO users (username, hashed_password, full_name, email, password_changed_at, created_at) 
VALUES('don01', NULL, 'don tech', 'don@gmail.com', now(), now());
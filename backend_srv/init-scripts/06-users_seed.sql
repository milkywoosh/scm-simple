
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
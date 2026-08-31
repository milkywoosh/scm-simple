
CREATE TABLE roles (
    id SERIAL PRIMARY KEY ,
    rolename VARCHAR(50) UNIQUE,
    "description" VARCHAR(90),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- 

CREATE TABLE user_roles (
    user_id bigint,
    role_id bigint
);

CREATE INDEX idx_userroles
    ON user_roles (user_id, role_id);

-- foreign key to [username, role]
ALTER TABLE user_roles
    ADD CONSTRAINT fk_users_on_userroles
        FOREIGN KEY (user_id)
        REFERENCES users(id);

ALTER TABLE user_roles
    ADD CONSTRAINT fk_roles_on_userroles
        FOREIGN KEY (role_id)
        REFERENCES roles(id);

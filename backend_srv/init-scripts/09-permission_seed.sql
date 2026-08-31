CREATE TABLE permissions (
    id PRIMARY KEY SERIAL,
    "name" VARCHAR(50) UNIQUE, -- 'procurement:create' ; 'procurement:delete'; 'procurement:approve' 
    "resource" VARCHAR(90), -- procurement
    "action" VARCHAR(90), -- create
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_by BIGINT, -- from user_id
);

CREATE TABLE role_permissions (
    id PRIMARY KEY SERIAL,
    role_id bigint,
    permission_id bigint
);

ALTER TABLE role_permissions
    ADD CONSTRAINT fk_roles_on_rolepermissions
        FOREIGN KEY (role_id)
        REFERENCES roles(id);

ALTER TABLE role_permissions
    ADD CONSTRAINT fk_permissions_on_rolepermissions
        FOREIGN KEY (permission_id)
        REFERENCES permissions(id);



-- 
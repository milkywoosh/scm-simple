CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    "name" VARCHAR(50) UNIQUE, -- 'procurement:create' ; 'procurement:delete'; 'procurement:approve' 
    "resource" VARCHAR(90), -- procurement
    "action" VARCHAR(90), -- create
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_by BIGINT -- from user_id
);

CREATE TABLE role_permissions (
    role_id bigint,
    permission_id bigint
);

-- on
ALTER TABLE role_permissions 
ADD CONSTRAINT unique_role_permissions_fields UNIQUE (role_id, permission_id);

ALTER TABLE role_permissions
    ADD CONSTRAINT fk_roles_on_rolepermissions
        FOREIGN KEY (role_id)
        REFERENCES roles(id);

ALTER TABLE role_permissions
    ADD CONSTRAINT fk_permissions_on_rolepermissions
        FOREIGN KEY (permission_id)
        REFERENCES permissions(id);



-- 
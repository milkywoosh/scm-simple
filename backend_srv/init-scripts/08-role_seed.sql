
CREATE TABLE roles (
    id SERIAL PRIMARY KEY ,
    rolename VARCHAR(50) UNIQUE,
    "description" VARCHAR(90),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

INSERT INTO roles (
    rolename,
    "description"
) VALUES 
    (
        'warehouse_manager', 'manager of warehouse'
    ),
    (
        'warehouse_spv', 'supervisor of warehouse'
    ),
    (
        'warehouse_staff', 'staff of warehouse'
    );

INSERT INTO roles (
    rolename,
    "description"
) VALUES 
    (
        'technician', 'person who is incharge for installation, dismantling, deployment'
    );
    

CREATE TABLE user_roles (
    user_id bigint,
    role_id bigint,
    created_at timestamptz default now()
);

CREATE INDEX idx_userroles ON user_roles (user_id, role_id);

-- on
ALTER TABLE user_roles 
ADD CONSTRAINT unique_user_roles_fields UNIQUE (user_id, role_id);


-- foreign key to [username, role]
ALTER TABLE user_roles
    ADD CONSTRAINT fk_users_on_userroles
        FOREIGN KEY (user_id)
        REFERENCES users(id);

ALTER TABLE user_roles
    ADD CONSTRAINT fk_roles_on_userroles
        FOREIGN KEY (role_id)
        REFERENCES roles(id);



INSERT INTO user_roles (user_id, role_id)
    SELECT u.id, r.id FROM users u, roles r
        WHERE u.username = 'ben01' AND r.rolename = 'warehouse_staff'
    UNION ALL
    SELECT u.id, r.id FROM users u, roles r
        WHERE u.username = 'ron01' AND r.rolename = 'warehouse_manager';

-- assign sample user for technician
INSERT INTO user_roles (user_id, role_id)
    SELECT u.id, r.id FROM users u, roles r
        WHERE u.username = 'don01' AND r.rolename = 'technician'
    





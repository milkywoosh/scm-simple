CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    "name" VARCHAR(50) UNIQUE, -- etc 'procurement:create' ; 'procurement:delete'; 'procurement:approve' 
    "resource" type_of_trans, -- procurement **defined type at transaction_seed.sql
    "action" VARCHAR(90), -- create
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_by BIGINT -- from user_id
);

-- **defined initially at transaction_seed.sql
ALTER TYPE type_of_trans ADD VALUE 'user_administration';

INSERT INTO permissions (
    "name", "resource", "action", created_by
) VALUES 
    (
        'warehouse_to_warehouse:create',
        'warehouse_to_warehouse',
        'create',
        999 -- for admin user 
    ),
    (
        'warehouse_to_warehouse:submit',
        'warehouse_to_warehouse',
        'submit',
        999 -- for admin user 
    ),
    (
        'warehouse_to_warehouse:reject',
        'warehouse_to_warehouse',
        'reject',
        999 -- for admin user 
    ),
    (
        'warehouse_to_warehouse:cancel',
        'warehouse_to_warehouse',
        'cancel',
        999 -- for admin user 
    ),
    (
        'warehouse_to_warehouse:approve',
        'warehouse_to_warehouse',
        'approve',
        999 -- for admin user 
    );

CREATE TABLE role_permissions (
    role_id bigint,
    permission_id bigint,
    created_at timestamptz default now(),
    created_by BIGINT NOT NULL
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


INSERT INTO role_permissions (
    role_id,
    permission_id,
    created_by
) SELECT r.id, p.id, 999 FROM roles r, permissions p 
    WHERE r.rolename = 'warehouse_manager' AND p."name" = 'warehouse_to_warehouse:cancel'
  UNION ALL
  SELECT r.id, p.id, 999 FROM roles r, permissions p 
    WHERE r.rolename = 'warehouse_manager' AND p."name" = 'warehouse_to_warehouse:approve';

-- ok
INSERT INTO role_permissions (
    role_id,
    permission_id,
    created_by
) SELECT r.id, p.id, 999 FROM roles r, permissions p 
    WHERE r.rolename = 'warehouse_staff' AND p."name" = 'warehouse_to_warehouse:create'
  UNION ALL
  SELECT r.id, p.id, 999 FROM roles r, permissions p 
    WHERE r.rolename = 'warehouse_staff' AND p."name" = 'warehouse_to_warehouse:submit'
  UNION ALL
  SELECT r.id, p.id, 999 FROM roles r, permissions p 
    WHERE r.rolename = 'warehouse_staff' AND p."name" = 'warehouse_to_warehouse:reject';

-- 
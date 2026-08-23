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


-- seed data
INSERT INTO warehouses (
    location_code,
    "description",
    parent_location_code
) VALUES (
    'PGC001',
    'pusat grosir cililitan',
    null -- last gakbole comma
);


INSERT INTO warehouses (
    location_code,
    "description",
    parent_location_code
) VALUES (
    'BOLUJAYA01',
    'Bolu Jaya, jl. lembur 001',
    'PGC001' -- last gakbole comma
);

INSERT INTO warehouses (
    location_code,
    "description",
    parent_location_code
) VALUES (
    'BOLUJAYA02',
    'Bolu Jaya, jl. nacim 002',
    'PGC001' -- last gakbole comma
);

INSERT INTO warehouses (
    location_code,
    "description",
    parent_location_code
) VALUES (
    'BOLUJAYA03',
    'Bolu Jaya, jl. cikeas 002',
    'PGC001'-- last gakbole comma
);
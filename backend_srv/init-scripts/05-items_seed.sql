
-- ENUM TYPES

CREATE TYPE item_status_enum AS ENUM (
    'ALLOCATED',
    'INTECH',
    'AVAILABLE',
    'INSTALLED',
    'DEPLOYED',
    'INBETWEEN',
    'DAMAGE',
    'INLOCKED'
);


-- ITEMS

CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    serial_number VARCHAR(15) NOT NULL UNIQUE,
    factory_serial_number VARCHAR(150),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    curr_status item_status_enum NOT NULL,
    curr_transaction VARCHAR(60),
    curr_location_code VARCHAR(50),
    product_code VARCHAR(90),
    introduction_number VARCHAR(90), -- awal mula munculnya item (bisa procurement, migration, pinjaman)
    CONSTRAINT fk_items_curr_transaction
        FOREIGN KEY (curr_transaction)
        REFERENCES transaction_item_transfers(transaction_number)
);


-- ITEM HISTORIES
CREATE TABLE item_histories (
    id BIGSERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL,
    transaction_number VARCHAR(60) NOT NULL,
    previous_status item_status_enum NOT NULL, -- for disallocating
    status_history item_status_enum NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_item_histories_item
        FOREIGN KEY (item_id)
        REFERENCES items(id)
);


CREATE INDEX idx_item_histories_item_transaction
    ON item_histories (item_id, transaction_number);



-- data seeder

INSERT INTO items (
    serial_number,
    factory_serial_number,
    created_at,
    curr_status,
    curr_transaction,
    curr_location_code,
    product_code,
    introduction_number
) VALUES
(
    'ITM000000000001',
    'FACTORY-SN-AX2026-00001',
    '2026-01-10 08:30:00+07',
    'AVAILABLE',
    'TRX-2026-000001',
    'WH-JKT-01',
    'PROD-MTR-001',
    'PROC-2026-00001'
),
(
    'ITM000000000002',
    'FACTORY-SN-AX2026-00002',
    '2026-01-12 09:15:00+07',
    'INTECH',
    'TRX-2026-000002',
    'TECH-JKT-01',
    'PROD-MTR-001',
    'PROC-2026-00001'
),
(
    'ITM000000000003',
    'FACTORY-SN-AX2026-00003',
    '2026-01-15 10:00:00+07',
    'ALLOCATED',
    'TRX-2026-000003',
    'WH-BDG-01',
    'PROD-MTR-002',
    'PROC-2026-00002'
),
(
    'ITM000000000004',
    'FACTORY-SN-AX2026-00004',
    '2026-01-20 11:30:00+07',
    'INSTALLED',
    'TRX-2026-000004',
    'SITE-BDG-001',
    'PROD-MTR-002',
    'PROC-2026-00002'
),
(
    'ITM000000000005',
    'FACTORY-SN-AX2026-00005',
    '2026-01-25 13:45:00+07',
    'DEPLOYED',
    'TRX-2026-000005',
    'SITE-JKT-001',
    'PROD-MTR-003',
    'MIG-2026-00001'
);


INSERT INTO item_histories (
    item_id,
    transaction_number,
    previous_status,
    status_history,
    created_at
)
SELECT
    i.id,
    x.transaction_number,
    x.previous_status::item_status_enum,
    x.status_history::item_status_enum,
    x.created_at
FROM (
    VALUES
        ('ITM000000000001', 'TRX-2026-000001', 'ALLOCATED', 'AVAILABLE', '2026-01-10 08:30:00+07'::timestamptz),
        ('ITM000000000002', 'TRX-2026-000002', 'ALLOCATED', 'AVAILABLE', '2026-01-12 09:15:00+07'::timestamptz),
        ('ITM000000000002', 'TRX-2026-000002', 'ALLOCATED', 'INTECH',    '2026-01-13 14:20:00+07'::timestamptz)
) AS x(
    serial_number,
    transaction_number,
    previous_status,
    status_history,
    created_at
)
JOIN items i
    ON i.serial_number = x.serial_number;
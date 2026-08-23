
CREATE TYPE type_of_trans AS ENUM ('warehouse_to_warehouse', 'warehouse_to_technician', 'technician_to_warehouse', 'procurement', 'installation', 'deployment');

CREATE TABLE transaction_item_transfers (
    id serial primary key ,
    transaction_number VARCHAR(60) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    transaction_type type_of_trans,
    origin VARCHAR(60),
    destination VARCHAR(60),
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    submitted_at DATE,
    approved_at DATE,
    canceled_at DATE
);

CREATE INDEX trans_number_idx ON transaction_item_transfers (transaction_number);

CREATE TABLE transaction_item_transfer_details (
    id  serial primary key ,
    id_trans_item_transfer INT,
    identifier_item VARCHAR(80),
    added_at DATE NOT NULL DEFAULT CURRENT_DATE
);
ALTER TABLE transaction_item_transfer_details
    ADD CONSTRAINT fk_transaction_item_transfer
        FOREIGN KEY (id_trans_item_transfer)
        REFERENCES transaction_item_transfers(id)
        ON DELETE CASCADE;



-- seed data

-- =========================================================
-- TRANSACTION ITEM TRANSFERS
-- =========================================================

INSERT INTO transaction_item_transfers (
    transaction_number,
    status,
    transaction_type,
    origin,
    destination,
    created_at,
    submitted_at,
    approved_at
) VALUES
(
    'TRX-2026-000001',
    'approved',
    'procurement',
    'FACTORY',
    'WH-JKT-01',
    '2026-01-10',
    '2026-01-10',
    '2026-01-10'
),
(
    'TRX-2026-000002',
    'approved',
    'warehouse_to_warehouse',
    'WH-JKT-01',
    'TECH-JKT-01',
    '2026-01-12',
    '2026-01-12',
    '2026-01-12'
),
(
    'TRX-2026-000003',
    'approved',
    'warehouse_to_warehouse',
    'WH-JKT-01',
    'WH-BDG-01',
    '2026-01-15',
    '2026-01-15',
    '2026-01-15'
),
(
    'TRX-2026-000004',
    'approved',
    'installation',
    'WH-BDG-01',
    'SITE-BDG-001',
    '2026-01-20',
    '2026-01-20',
    '2026-01-20'
),
(
    'TRX-2026-000005',
    'approved',
    'deployment',
    'SITE-BDG-001',
    'SITE-JKT-001',
    '2026-01-25',
    '2026-01-25',
    '2026-01-25'
);


-- =========================================================
-- TRANSACTION ITEM TRANSFER DETAILS
-- =========================================================

INSERT INTO transaction_item_transfer_details (
    id_trans_item_transfer,
    identifier_item,
    added_at
) VALUES
(
    (SELECT id
     FROM transaction_item_transfers
     WHERE transaction_number = 'TRX-2026-000001'),
    'ITM000000000001',
    '2026-01-10'
),
(
    (SELECT id
     FROM transaction_item_transfers
     WHERE transaction_number = 'TRX-2026-000002'),
    'ITM000000000002',
    '2026-01-12'
),
(
    (SELECT id
     FROM transaction_item_transfers
     WHERE transaction_number = 'TRX-2026-000003'),
    'ITM000000000003',
    '2026-01-15'
),
(
    (SELECT id
     FROM transaction_item_transfers
     WHERE transaction_number = 'TRX-2026-000004'),
    'ITM000000000004',
    '2026-01-20'
),
(
    (SELECT id
     FROM transaction_item_transfers
     WHERE transaction_number = 'TRX-2026-000005'),
    'ITM000000000005',
    '2026-01-25'
);



-- table transaction_transfer_tracker => untuk tracking proses antara send dan receive
-- id serial primary key
-- outbound_transaction string (fk from transaction_item_transfers)
-- status_outbound string (enum type [ondelivery, arrived])
-- inbound_number string (fk from transaction_item_transfers)
-- status_inbound string (enum type [di])

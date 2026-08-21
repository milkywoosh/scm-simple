
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
    proc_number VARCHAR(90),
    CONSTRAINT fk_items_curr_transaction
        FOREIGN KEY (curr_transaction)
        REFERENCES transaction_item_transfers(transaction_number)
);


-- ITEM HISTORIES
CREATE TABLE item_histories (
    id BIGSERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL,
    transaction_number VARCHAR(60) NOT NULL,
    status_history item_status_enum NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_item_histories_item
        FOREIGN KEY (item_id)
        REFERENCES items(id)
);


CREATE INDEX idx_item_histories_item_transaction
    ON item_histories (item_id, transaction_number);
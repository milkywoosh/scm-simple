
CREATE TYPE type_of_trans AS ENUM ('warehouse_to_warehouse', 'warehouse_to_technician', 'technician_to_warehouse');

CREATE TABLE transaction_item_transfers (
    id serial primary key ,
    transaction_number VARCHAR(60) NOT NULL UNIQUE,
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
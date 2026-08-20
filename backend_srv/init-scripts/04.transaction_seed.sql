
CREATE TABLE transaction_item_transfers (
    id number  primary key serial,
    transaction_number VARCHAR(60) NOT NULL UNIQUE,
    origin VARCHAR(60),
    destination VARCHAR(60),
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    submitted_at DATE,
    approved_at DATE,
    canceled_at DATE,
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
-- table transaction_transfer_tracker => untuk tracking proses antara send dan receive

CREATE TABLE transaction_transfer_tracker (
    id serial primary key,
    outbound_transaction VARCHAR(60), -- delete if transaction_number is canceled
    inbound_transaction VARCHAR(60),
    delivery_at TIMESTAMPTZ DEFAULT NULL, -- jadi status inbound outbound, parsing jadi status
    arrived_at TIMESTAMPTZ DEFAULT NULL,-- jadi status inbound outbound, parsing jadi status
    files BIGINT-- tracker_file_id -- add later
);

ALTER TABLE transaction_transfer_tracker
    ADD CONSTRAINT fk_outbound_transaction
        FOREIGN KEY (outbound_transaction)
        REFERENCES transaction_item_transfers(transaction_number);

ALTER TABLE transaction_transfer_tracker
    ADD CONSTRAINT fk_inbound_transaction
        FOREIGN KEY (inbound_transaction)
        REFERENCES transaction_item_transfers(transaction_number);

CREATE INDEX tracker_outbound_idx ON transaction_transfer_tracker (outbound_transaction);
CREATE INDEX tracker_inbound_idx ON transaction_transfer_tracker (inbound_transaction);
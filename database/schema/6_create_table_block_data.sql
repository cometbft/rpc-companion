-- Table: comet.block_data

CREATE TABLE IF NOT EXISTS comet.block_data
(
    height bigint NOT NULL,
    transaction bytea NOT NULL,
    CONSTRAINT block_height_fk FOREIGN KEY (height)
        REFERENCES comet.block_result (block_header_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
) TABLESPACE pg_default;

-- Table: comet.last_commit_signature

CREATE TABLE IF NOT EXISTS comet.last_commit_signature
(
    height bigint NOT NULL,
    block_id_flag comet.uint8 NOT NULL,
    validator_address bytea NOT NULL,
    signature_timestamp timestamp with time zone NOT NULL,
    signature bytea NOT NULL,
    CONSTRAINT last_commit_height_fk FOREIGN KEY (height)
        REFERENCES comet.result_block (block_last_commit_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
) TABLESPACE pg_default;
-- Table: comet.block_commit_sig

DROP TABLE IF EXISTS comet.block_commit_sig;
CREATE TABLE comet.block_commit_sig
(
    height bigint NOT NULL,
    block_id_flag comet.uint8 NOT NULL,
    validator_address bytea,
    timestamp timestamp with time zone NOT NULL,
    signature bytea,
    CONSTRAINT block_commit_sig_height_fk FOREIGN KEY (height)
        REFERENCES comet.block_result (block_last_commit_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
) TABLESPACE pg_default;
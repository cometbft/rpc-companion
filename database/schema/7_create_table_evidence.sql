-- Table: comet.evidence

CREATE TABLE IF NOT EXISTS comet.evidence
(
    height bigint NOT NULL,
    evidence_type bytea NOT NULL,
    vote_a_type int NOT NULL,
    vote_b_type int NOT NULL,
    CONSTRAINT evidence_height_fk FOREIGN KEY (height)
        REFERENCES comet.result_block (block_header_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
) TABLESPACE pg_default;
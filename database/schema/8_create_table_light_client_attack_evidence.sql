-- Table: comet.light_client_attack_evidence

CREATE TABLE IF NOT EXISTS comet.light_client_attack_evidence
(
    height bigint NOT NULL,
    common_height bigint NOT NULL,
    CONSTRAINT evidence_height_fk FOREIGN KEY (height)
        REFERENCES comet.result_block (block_header_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
) TABLESPACE pg_default;
-- Table: comet.light_client_attack_evidence
DROP TABLE IF EXISTS comet.light_client_attack_evidence;
CREATE TABLE comet.light_client_attack_evidence
(
    height bigint NOT NULL,
    evidence_type bytea NOT NULL,
    common_height bigint NOT NULL,
    total_voting_power bigint NOT NULL,
    evidence_timestamp timestamp with time zone NOT NULL,
    CONSTRAINT lca_evidence_height_fk FOREIGN KEY (height)
        REFERENCES comet.result_block (block_header_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
) TABLESPACE pg_default;
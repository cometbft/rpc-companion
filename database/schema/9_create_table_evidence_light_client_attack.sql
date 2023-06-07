-- Table: comet.evidence_light_client_attack

DROP TABLE IF EXISTS comet.evidence_light_client_attack;
CREATE TABLE comet.evidence_light_client_attack
(
    height bigint NOT NULL,
    evidence_type bytea NOT NULL,
    common_height bigint NOT NULL,
    total_voting_power bigint NOT NULL,
    timestamp timestamp with time zone NOT NULL,
    CONSTRAINT lca_evidence_height_fk FOREIGN KEY (height)
        REFERENCES comet.block_result (block_header_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
) TABLESPACE pg_default;
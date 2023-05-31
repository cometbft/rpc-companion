-- Table: comet.duplicate_vote_evidence

CREATE TABLE IF NOT EXISTS comet.duplicate_vote_evidence
(
    height bigint NOT NULL,
    evidence_type bytea NOT NULL,

    vote_a_type int NOT NULL,
    vote_a_height bigint NOT NULL,
    vote_a_round int NOT NULL,
    vote_a_block_id_hash bytea NOT NULL,
    vote_a_block_id_parts_hash bytea NOT NULL,
    vote_a_block_id_parts_total comet.uint32 NOT NULL,
    vote_a_timestamp timestamp with time zone NOT NULL,
    vote_a_validator_address bytea NOT NULL,
    vote_a_validator_index int NOT NULL,
    vote_a_signature bytea NOT NULL,
    vote_a_extension bytea NOT NULL,
    vote_a_extension_signature bytea NOT NULL,

    vote_b_type int NOT NULL,
    vote_b_height bigint NOT NULL,
    vote_b_round int NOT NULL,
    vote_b_block_id_hash bytea NOT NULL,
    vote_b_block_id_parts_hash bytea NOT NULL,
    vote_b_block_id_parts_total comet.uint32 NOT NULL,
    vote_b_timestamp timestamp with time zone NOT NULL,
    vote_b_validator_address bytea NOT NULL,
    vote_b_validator_index int NOT NULL,
    vote_b_signature bytea NOT NULL,
    vote_b_extension bytea NOT NULL,
    vote_b_extension_signature bytea NOT NULL,

    total_voting_power bigint NOT NULL,
    validator_voting_power bigint NOT NULL,
    evidence_timestamp timestamp with time zone NOT NULL,
    CONSTRAINT evidence_height_fk FOREIGN KEY (height)
        REFERENCES comet.result_block (block_header_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
) TABLESPACE pg_default;
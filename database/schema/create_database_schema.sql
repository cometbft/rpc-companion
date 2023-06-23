-- SCHEMA: comet

CREATE SCHEMA IF NOT EXISTS comet;

-- DOMAIN: comet.uint8

DO $$ BEGIN
    CREATE DOMAIN comet.uint8 AS numeric;

    ALTER DOMAIN comet.uint8
        ADD CONSTRAINT value_max CHECK (VALUE <= '255'::numeric);

    ALTER DOMAIN comet.uint8
        ADD CONSTRAINT value_positive CHECK (VALUE >= 0::numeric);
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- DOMAIN: comet.uint32

DO $$ BEGIN
    CREATE DOMAIN comet.uint32 AS numeric;

    ALTER DOMAIN comet.uint32
        ADD CONSTRAINT value_max CHECK (VALUE <= '4294967295'::numeric);

    ALTER DOMAIN comet.uint32
        ADD CONSTRAINT value_positive CHECK (VALUE >= 0::numeric);
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- DOMAIN: comet.uint64

DO $$ BEGIN
    CREATE DOMAIN comet.uint64
        AS numeric;

    ALTER DOMAIN comet.uint64 OWNER TO postgres;

    ALTER DOMAIN comet.uint64
        ADD CONSTRAINT value_max CHECK (VALUE <= '18446744073709551615'::numeric);

    ALTER DOMAIN comet.uint64
        ADD CONSTRAINT value_positive CHECK (VALUE >= 0::numeric);
    EXCEPTION
        WHEN duplicate_object THEN null;
END $$;

-- TABLE: comet.block_result

DROP TABLE IF EXISTS comet.block_result CASCADE;

CREATE TABLE comet.block_result
(
    block_id_hash bytea NOT NULL,
    block_id_parts_hash bytea NOT NULL,
    block_id_parts_total comet.uint32 NOT NULL,
    block_header_height bigint NOT NULL,
    block_header_block_time timestamp with time zone NOT NULL,
    block_header_chain_id text NOT NULL,
    block_header_version_block comet.uint64 NOT NULL,
    block_header_version_app comet.uint64 NOT NULL,
    block_header_data_hash bytea NOT NULL,
    block_header_last_commit_hash bytea NOT NULL,
    block_header_validators_hash bytea NOT NULL,
    block_header_next_validators_hash bytea NOT NULL,
    block_header_consensus_hash bytea NOT NULL,
    block_header_app_hash bytea NOT NULL,
    block_header_last_results_hash bytea NOT NULL,
    block_header_evidence_hash bytea NOT NULL,
    block_header_proposer_address bytea NOT NULL,
    block_header_last_block_id_hash bytea NOT NULL,
    block_header_last_block_id_parts_hash bytea NOT NULL,
    block_header_last_block_id_part_total comet.uint32 NOT NULL,
    block_last_commit_height comet.uint64 NOT NULL,
    block_last_commit_round comet.uint32 NOT NULL,
    block_last_commit_block_id_hash bytea NOT NULL,
    block_last_commit_block_id_parts_total comet.uint32 NOT NULL,
    block_last_commit_block_id_parts_hash bytea NOT NULL,
    CONSTRAINT block_pkey PRIMARY KEY (block_header_height),
    CONSTRAINT last_commit_height_unique UNIQUE (block_last_commit_height),
    CONSTRAINT height_positive CHECK (block_header_height >= 0)
);

-- TABLE: comet.block_data

DROP TABLE IF EXISTS comet.block_data CASCADE;

CREATE TABLE comet.block_data
(
    height bigint NOT NULL,
    transaction bytea NOT NULL,
    CONSTRAINT block_height_fk FOREIGN KEY (height)
        REFERENCES comet.block_result (block_header_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- TABLE: comet.block_commit_sig

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
);

-- TABLE: comet.evidence_duplicate_vote

DROP TABLE IF EXISTS comet.evidence_duplicate_vote;

CREATE TABLE comet.evidence_duplicate_vote
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

    total_voting_power bigint NOT NULL,
    validator_voting_power bigint NOT NULL,
    evidence_timestamp timestamp with time zone NOT NULL,
    CONSTRAINT dv_evidence_height_fk FOREIGN KEY (height)
        REFERENCES comet.block_result (block_header_height) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- TABLE: comet.evidence_light_client_attack

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
);

DROP TABLE IF EXISTS comet.validator;

CREATE TABLE comet.validator
(
    id bigserial NOT NULL,
    address bytea NOT NULL,
    pub_key_type bytea NOT NULL,
    pub_key_value bytea NOT NULL,
    voting_power bigint NOT NULL,
    proposer_priority bigint NOT NULL,
    CONSTRAINT validator_pkey PRIMARY KEY (id),
    CONSTRAINT unique_validator UNIQUE (pub_key_type, pub_key_value, voting_power)
);

-- Table: comet.result_block

CREATE TABLE IF NOT EXISTS comet.result_block
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
    CONSTRAINT height_positive CHECK (block_header_height >= 0)
) TABLESPACE pg_default;

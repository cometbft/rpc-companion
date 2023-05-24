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
    block_last_block_id_hash bytea NOT NULL,
    block_last_block_id_parts_hash bytea NOT NULL,
    block_last_block_id_part_total comet.uint32 NOT NULL,
    block_data_hash bytea NOT NULL,
    CONSTRAINT block_pkey PRIMARY KEY (block_header_height),
    CONSTRAINT height_positive CHECK (block_header_height >= 0)
) TABLESPACE pg_default;

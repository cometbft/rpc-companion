-- Table: comet.block

CREATE TABLE IF NOT EXISTS comet.block
(
    height bigint NOT NULL,
    block_time timestamp with time zone NOT NULL,
    chain_id text NOT NULL,
    version_block comet.uint64 NOT NULL,
    version_app comet.uint64 NOT NULL,
    CONSTRAINT block_pkey PRIMARY KEY (height),
    CONSTRAINT height_positive CHECK (height >= 0)
) TABLESPACE pg_default;

-- Table: comet.block

CREATE TABLE IF NOT EXISTS comet.block
(
    height bigint NOT NULL,
    chain_id smallint NOT NULL,
    version integer NOT NULL,
    block_time timestamp with time zone NOT NULL,
    CONSTRAINT block_pkey PRIMARY KEY (height),
    CONSTRAINT chain_id_fk FOREIGN KEY (chain_id)
    REFERENCES comet.chain_id (chain_id) MATCH SIMPLE
    ON UPDATE RESTRICT
    ON DELETE RESTRICT,
    CONSTRAINT height_positive CHECK (height >= 0)
) TABLESPACE pg_default;

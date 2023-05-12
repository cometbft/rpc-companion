-- Table: comet.block

CREATE TABLE IF NOT EXISTS comet.block
(
    header_height bigint NOT NULL,
    header_chain_id smallint NOT NULL,
    header_version smallint NOT NULL,
    header_time timestamp with time zone NOT NULL,
    CONSTRAINT block_pkey PRIMARY KEY (header_height),
    CONSTRAINT chain_id_fk FOREIGN KEY (header_chain_id)
    REFERENCES comet.chain_id (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID,
    CONSTRAINT height_positive CHECK (header_height >= 0)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS comet.block
    OWNER to postgres;
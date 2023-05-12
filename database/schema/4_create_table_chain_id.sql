-- Table: comet.chain_id

CREATE TABLE IF NOT EXISTS comet.chain_id
(
    id smallserial NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT chain_id_pkey PRIMARY KEY (id),
    CONSTRAINT chain_id_name_key UNIQUE (name)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS comet.chain_id
    OWNER to postgres;
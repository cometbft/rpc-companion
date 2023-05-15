-- Table: comet.chain_id

CREATE TABLE IF NOT EXISTS comet.chain_id
(
    chain_id smallserial NOT NULL,
    chain_name text NOT NULL,
    CONSTRAINT chain_id_pkey PRIMARY KEY (chain_id),
    CONSTRAINT chain_id_name_key UNIQUE (chain_name)
)
TABLESPACE pg_default;
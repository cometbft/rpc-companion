-- Table: comet.consensus

CREATE TABLE IF NOT EXISTS comet.consensus
(
    id smallint NOT NULL DEFAULT nextval('consensus_id_seq'::regclass),
    block uint64 NOT NULL,
    app uint64 NOT NULL,
    CONSTRAINT consensus_pk PRIMARY KEY (block, app),
    CONSTRAINT block_app_unique UNIQUE (block, app)
    INCLUDE(block, app)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS comet.consensus
    OWNER to postgres;
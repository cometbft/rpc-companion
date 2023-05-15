-- Table: comet.consensus

CREATE TABLE IF NOT EXISTS comet.consensus
(
    consensus_id serial NOT NULL,
    block comet.uint64 NOT NULL,
    app comet.uint64 NOT NULL,
    CONSTRAINT consensus_pk PRIMARY KEY (block, app),
    CONSTRAINT block_app_unique UNIQUE (block, app)
    INCLUDE(block, app)
) TABLESPACE pg_default;

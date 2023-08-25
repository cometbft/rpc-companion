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

-- TABLE: comet.block

DROP TABLE IF EXISTS comet.block CASCADE;

CREATE TABLE comet.block
(
    height  comet.uint64 NOT NULL,
    data    bytea NOT NULL,
    CONSTRAINT block_pkey PRIMARY KEY (height)
);

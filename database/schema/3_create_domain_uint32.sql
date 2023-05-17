-- DOMAIN: comet.uint32

CREATE DOMAIN comet.uint32
    AS numeric;

ALTER DOMAIN comet.uint32
    ADD CONSTRAINT value_max CHECK (VALUE <= '4294967295'::numeric);

ALTER DOMAIN comet.uint32
    ADD CONSTRAINT value_positive CHECK (VALUE >= 0::numeric);
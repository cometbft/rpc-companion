-- DOMAIN: comet.uint8

CREATE DOMAIN comet.uint8
    AS numeric;

ALTER DOMAIN comet.uint8
    ADD CONSTRAINT value_max CHECK (VALUE <= '255'::numeric);

ALTER DOMAIN comet.uint8
    ADD CONSTRAINT value_positive CHECK (VALUE >= 0::numeric);
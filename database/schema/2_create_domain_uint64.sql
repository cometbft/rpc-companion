-- DOMAIN: comet.uint64

CREATE DOMAIN comet.uint64
    AS numeric;

ALTER DOMAIN comet.uint64
    ADD CONSTRAINT value_max CHECK (VALUE <= '18446744073709551615'::numeric);

ALTER DOMAIN comet.uint64
    ADD CONSTRAINT value_positive CHECK (VALUE >= 0::numeric);
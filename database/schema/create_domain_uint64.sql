-- DOMAIN: comet.uint64

CREATE DOMAIN public.uint64
    AS numeric;

ALTER DOMAIN public.uint64 OWNER TO postgres;

ALTER DOMAIN public.uint64
    ADD CONSTRAINT value_max CHECK (VALUE <= '18446744073709551615'::numeric);

ALTER DOMAIN public.uint64
    ADD CONSTRAINT value_positive CHECK (VALUE >= 0::numeric);
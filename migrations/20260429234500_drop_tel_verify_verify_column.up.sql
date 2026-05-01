SET statement_timeout = 0;

--bun:split
ALTER TABLE tel_verify
DROP COLUMN IF EXISTS verify;

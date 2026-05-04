SET statement_timeout = 0;

--bun:split
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;

--bun:split
ALTER TABLE tel_verify DROP COLUMN IF EXISTS signup_first_name;

--bun:split
ALTER TABLE tel_verify DROP COLUMN IF EXISTS signup_last_name;

--bun:split
ALTER TABLE tel_verify DROP COLUMN IF EXISTS signup_email;

--bun:split
ALTER TABLE tel_verify DROP COLUMN IF EXISTS password_hash;

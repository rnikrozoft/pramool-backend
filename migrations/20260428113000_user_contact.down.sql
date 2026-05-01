SET statement_timeout = 0;

--bun:split
ALTER TABLE users
DROP COLUMN IF EXISTS facebook,
DROP COLUMN IF EXISTS email;

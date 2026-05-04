SET statement_timeout = 0;

--bun:split
ALTER TABLE tel_verify
ADD COLUMN IF NOT EXISTS signup_first_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS signup_last_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS signup_email VARCHAR(120) DEFAULT '',
ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);

--bun:split
ALTER TABLE users
ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);

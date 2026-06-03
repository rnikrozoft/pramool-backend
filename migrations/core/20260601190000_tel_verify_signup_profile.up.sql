SET statement_timeout = 0;

--bun:split

-- Idempotent: signup step stores profile on tel_verify until POST /users completes onboarding.
ALTER TABLE tel_verify ADD COLUMN IF NOT EXISTS signup_first_name VARCHAR(100);
ALTER TABLE tel_verify ADD COLUMN IF NOT EXISTS signup_last_name VARCHAR(100);
ALTER TABLE tel_verify ADD COLUMN IF NOT EXISTS signup_email VARCHAR(120) NOT NULL DEFAULT '';

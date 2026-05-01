SET statement_timeout = 0;

--bun:split
ALTER TABLE tel_verify
DROP COLUMN IF EXISTS otp_timeout_count,
DROP COLUMN IF EXISTS otp_banned_until;

SET statement_timeout = 0;

--bun:split
ALTER TABLE tel_verify
ADD COLUMN IF NOT EXISTS otp_timeout_count INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS otp_banned_until TIMESTAMPTZ;

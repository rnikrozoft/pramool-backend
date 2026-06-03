SET statement_timeout = 0;

--bun:split

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS qr_code_url TEXT;

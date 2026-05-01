SET statement_timeout = 0;

--bun:split
ALTER TABLE tel_verify
ADD COLUMN IF NOT EXISTS verify BOOLEAN DEFAULT FALSE;

SET statement_timeout = 0;

--bun:split
ALTER TABLE auctions
ADD COLUMN IF NOT EXISTS allow_early_close BOOLEAN NOT NULL DEFAULT FALSE;

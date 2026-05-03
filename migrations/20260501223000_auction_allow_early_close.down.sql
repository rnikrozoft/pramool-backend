SET statement_timeout = 0;

--bun:split
ALTER TABLE auctions
DROP COLUMN IF EXISTS allow_early_close;

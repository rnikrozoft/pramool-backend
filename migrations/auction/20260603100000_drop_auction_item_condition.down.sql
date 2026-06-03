SET statement_timeout = 0;

--bun:split

ALTER TABLE auctions ADD COLUMN IF NOT EXISTS item_condition VARCHAR(100) NOT NULL DEFAULT '';

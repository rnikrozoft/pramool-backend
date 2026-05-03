SET statement_timeout = 0;

--bun:split
ALTER TABLE auctions ADD COLUMN IF NOT EXISTS buy_now_price BIGINT NOT NULL DEFAULT 0;

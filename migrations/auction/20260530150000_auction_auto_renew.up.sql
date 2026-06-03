SET statement_timeout = 0;

--bun:split

ALTER TABLE auctions
    ADD COLUMN IF NOT EXISTS auto_renew BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS auto_renew_count BIGINT NOT NULL DEFAULT 0;

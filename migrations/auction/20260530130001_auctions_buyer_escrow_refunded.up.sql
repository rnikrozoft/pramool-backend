SET statement_timeout = 0;

--bun:split

ALTER TABLE auctions
    ADD COLUMN IF NOT EXISTS buyer_escrow_refunded_at TIMESTAMPTZ;

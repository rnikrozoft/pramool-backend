SET statement_timeout = 0;

--bun:split

ALTER TABLE auction_seller_reviews
    ADD COLUMN IF NOT EXISTS comment TEXT NOT NULL DEFAULT '';

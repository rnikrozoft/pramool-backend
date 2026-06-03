SET statement_timeout = 0;

--bun:split

ALTER TABLE auction_seller_reviews
    DROP COLUMN IF EXISTS comment;

SET statement_timeout = 0;

--bun:split

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS seller_review_points_total BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS buyer_reputation_points BIGINT NOT NULL DEFAULT 0;

--bun:split

UPDATE users
SET seller_review_points_total = COALESCE(reputation_points, 0),
    buyer_reputation_points = 0;

--bun:split

ALTER TABLE users
    DROP COLUMN IF EXISTS reputation_points,
    DROP COLUMN IF EXISTS seller_review_rating_points_total;

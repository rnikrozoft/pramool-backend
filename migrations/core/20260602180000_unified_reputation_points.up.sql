SET statement_timeout = 0;

--bun:split

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS reputation_points BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS seller_review_rating_points_total BIGINT NOT NULL DEFAULT 0;

--bun:split

UPDATE users u
SET seller_review_rating_points_total = COALESCE((
    SELECT SUM(r.seller_points)::bigint
    FROM auction_seller_reviews r
    WHERE r.seller_id = u.user_id
), 0);

--bun:split

UPDATE users
SET reputation_points = COALESCE(seller_review_points_total, 0) + COALESCE(buyer_reputation_points, 0);

--bun:split

ALTER TABLE users
    DROP COLUMN IF EXISTS seller_review_points_total,
    DROP COLUMN IF EXISTS buyer_reputation_points;

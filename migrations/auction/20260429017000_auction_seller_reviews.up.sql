SET statement_timeout = 0;

--bun:split

CREATE TABLE auction_seller_reviews (
    auction_id VARCHAR(40) PRIMARY KEY REFERENCES auctions(auction_id) ON DELETE CASCADE,
    buyer_user_id VARCHAR(13) NOT NULL,
    seller_id VARCHAR(13) NOT NULL,
    rating NUMERIC(3, 1) NOT NULL,
    seller_points INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT auction_seller_reviews_rating_range CHECK (rating >= 0.5 AND rating <= 5.0),
    CONSTRAINT auction_seller_reviews_points_positive CHECK (seller_points >= 1 AND seller_points <= 10)
);

--bun:split

CREATE INDEX idx_auction_seller_reviews_seller_created
    ON auction_seller_reviews(seller_id, created_at DESC);

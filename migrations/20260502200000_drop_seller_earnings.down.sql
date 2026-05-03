CREATE TABLE IF NOT EXISTS seller_earnings (
    earning_id BIGSERIAL PRIMARY KEY,
    seller_id VARCHAR(13) NOT NULL,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    winner_user_id VARCHAR(13) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'settled',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (auction_id)
);

CREATE INDEX IF NOT EXISTS idx_seller_earnings_seller_created_at ON seller_earnings(seller_id, created_at DESC);

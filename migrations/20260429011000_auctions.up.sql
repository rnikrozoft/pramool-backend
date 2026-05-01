SET statement_timeout = 0;

--bun:split
CREATE TABLE IF NOT EXISTS auctions (
    auction_id VARCHAR(40) PRIMARY KEY,
    seller_id VARCHAR(13) NOT NULL,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    item_condition VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_price BIGINT NOT NULL,
    bid_step BIGINT NOT NULL,
    current_bid BIGINT NOT NULL DEFAULT 0,
    total_bids BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'active',
    end_at TIMESTAMPTZ NOT NULL,
    cover_image_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--bun:split
CREATE TABLE IF NOT EXISTS auction_images (
    id BIGSERIAL PRIMARY KEY,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

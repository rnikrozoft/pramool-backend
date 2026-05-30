SET statement_timeout = 0;

--bun:split

CREATE TABLE auction_images (
    id BIGSERIAL PRIMARY KEY,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

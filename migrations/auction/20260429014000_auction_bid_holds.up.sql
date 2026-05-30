SET statement_timeout = 0;

--bun:split

CREATE TABLE auction_bid_holds (
    hold_id BIGSERIAL PRIMARY KEY,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    user_id VARCHAR(13) NOT NULL,
    held_amount BIGINT NOT NULL,
    hold_status VARCHAR(20) NOT NULL DEFAULT 'held',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    released_at TIMESTAMPTZ,
    settled_at TIMESTAMPTZ,
    UNIQUE (auction_id, user_id)
);

--bun:split

CREATE INDEX idx_auction_bid_holds_auction_status ON auction_bid_holds(auction_id, hold_status);

--bun:split

CREATE INDEX idx_auction_bid_holds_user_status ON auction_bid_holds(user_id, hold_status);

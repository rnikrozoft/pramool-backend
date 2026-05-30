SET statement_timeout = 0;

--bun:split

CREATE TABLE auction_bids (
    bid_id BIGSERIAL PRIMARY KEY,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    bidder_user_id VARCHAR(13) NOT NULL,
    bid_amount BIGINT NOT NULL,
    placed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (auction_id, bidder_user_id)
);

--bun:split

CREATE INDEX idx_auction_bids_auction_placed_at ON auction_bids(auction_id, placed_at DESC);

--bun:split

CREATE INDEX idx_auction_bids_bidder_placed_at ON auction_bids(bidder_user_id, placed_at DESC);

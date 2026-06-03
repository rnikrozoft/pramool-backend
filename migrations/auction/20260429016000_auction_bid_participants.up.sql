SET statement_timeout = 0;

--bun:split

CREATE TABLE auction_bid_participants (
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    bidder_user_id UUID NOT NULL,
    max_bid_amount BIGINT NOT NULL,
    last_bid_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (auction_id, bidder_user_id)
);

--bun:split

CREATE INDEX idx_auction_bid_participants_bidder_last
    ON auction_bid_participants(bidder_user_id, last_bid_at DESC);

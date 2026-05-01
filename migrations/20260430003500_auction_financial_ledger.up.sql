SET statement_timeout = 0;

--bun:split
ALTER TABLE auctions
ADD COLUMN IF NOT EXISTS winner_id VARCHAR(13),
ADD COLUMN IF NOT EXISTS settled_at TIMESTAMPTZ;

--bun:split
CREATE TABLE IF NOT EXISTS auction_bids (
    bid_id BIGSERIAL PRIMARY KEY,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    bidder_user_id VARCHAR(13) NOT NULL,
    bid_amount BIGINT NOT NULL,
    placed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--bun:split
CREATE INDEX IF NOT EXISTS idx_auction_bids_auction_placed_at ON auction_bids(auction_id, placed_at DESC);

--bun:split
CREATE INDEX IF NOT EXISTS idx_auction_bids_bidder_placed_at ON auction_bids(bidder_user_id, placed_at DESC);

--bun:split
CREATE TABLE IF NOT EXISTS auction_bid_holds (
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
CREATE INDEX IF NOT EXISTS idx_auction_bid_holds_auction_status ON auction_bid_holds(auction_id, hold_status);

--bun:split
CREATE INDEX IF NOT EXISTS idx_auction_bid_holds_user_status ON auction_bid_holds(user_id, hold_status);

--bun:split
CREATE TABLE IF NOT EXISTS bid_transactions (
    bid_tx_id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(13) NOT NULL,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    tx_type VARCHAR(30) NOT NULL,
    amount BIGINT NOT NULL,
    balance_before BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--bun:split
CREATE INDEX IF NOT EXISTS idx_bid_transactions_user_created_at ON bid_transactions(user_id, created_at DESC);

--bun:split
CREATE INDEX IF NOT EXISTS idx_bid_transactions_auction_created_at ON bid_transactions(auction_id, created_at DESC);

--bun:split
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

--bun:split
CREATE INDEX IF NOT EXISTS idx_seller_earnings_seller_created_at ON seller_earnings(seller_id, created_at DESC);

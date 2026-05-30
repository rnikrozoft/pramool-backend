SET statement_timeout = 0;

--bun:split

CREATE TABLE bid_transactions (
    bid_tx_id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(13) NOT NULL,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    tx_type VARCHAR(30) NOT NULL,
    amount BIGINT NOT NULL,
    bid_amount BIGINT,
    balance_before BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--bun:split

CREATE INDEX idx_bid_transactions_user_created_at ON bid_transactions(user_id, created_at DESC);

--bun:split

CREATE INDEX idx_bid_transactions_auction_created_at ON bid_transactions(auction_id, created_at DESC);

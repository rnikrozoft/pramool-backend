SET statement_timeout = 0;

--bun:split

CREATE TABLE auctions (
    auction_id VARCHAR(40) PRIMARY KEY,
    seller_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    category TEXT NOT NULL,
    item_condition VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_price BIGINT NOT NULL,
    bid_step BIGINT NOT NULL,
    current_bid BIGINT NOT NULL DEFAULT 0,
    total_bids BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'active',
    end_at TIMESTAMPTZ NOT NULL,
    allow_early_close BOOLEAN NOT NULL DEFAULT FALSE,
    early_close_hold_amount BIGINT NOT NULL DEFAULT 0,
    buy_now_price BIGINT NOT NULL DEFAULT 0,
    cover_image_url TEXT NOT NULL,
    winner_id UUID,
    settled_at TIMESTAMPTZ,
    seller_shipped_at TIMESTAMPTZ,
    buyer_received_at TIMESTAMPTZ,
    seller_payout_at TIMESTAMPTZ,
    payout_early_close BOOLEAN NOT NULL DEFAULT FALSE,
    seller_close_pause_bids_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

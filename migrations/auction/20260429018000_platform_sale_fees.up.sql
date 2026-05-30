SET statement_timeout = 0;

--bun:split

CREATE TABLE platform_sale_fees (
    platform_sale_fee_id BIGSERIAL PRIMARY KEY,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions(auction_id) ON DELETE CASCADE,
    seller_id VARCHAR(13) NOT NULL,
    winner_user_id VARCHAR(13) NOT NULL,
    winner_escrow_amount BIGINT NOT NULL CHECK (winner_escrow_amount > 0),
    seller_share_amount BIGINT NOT NULL CHECK (seller_share_amount >= 0),
    platform_fee_amount BIGINT NOT NULL CHECK (platform_fee_amount >= 0),
    payout_early_close BOOLEAN NOT NULL DEFAULT FALSE,
    seller_keep_pct SMALLINT NOT NULL CHECK (seller_keep_pct >= 0 AND seller_keep_pct <= 100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (auction_id),
    CHECK (seller_share_amount + platform_fee_amount = winner_escrow_amount)
);

--bun:split

CREATE INDEX idx_platform_sale_fees_created_at ON platform_sale_fees(created_at DESC);

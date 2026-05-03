SET statement_timeout = 0;

ALTER TABLE auctions
    DROP COLUMN IF EXISTS seller_shipped_at,
    DROP COLUMN IF EXISTS buyer_received_at,
    DROP COLUMN IF EXISTS seller_payout_at,
    DROP COLUMN IF EXISTS payout_early_close;

-- Deposit held from seller when allow_early_close is enabled (refunded if auction ends normally).
ALTER TABLE auctions
ADD COLUMN IF NOT EXISTS early_close_hold_amount BIGINT NOT NULL DEFAULT 0;

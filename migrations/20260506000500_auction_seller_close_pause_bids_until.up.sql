-- Pause window used by seller early-close flow to block incoming bids briefly before settlement.
ALTER TABLE auctions
ADD COLUMN IF NOT EXISTS seller_close_pause_bids_until TIMESTAMPTZ NULL;

SET statement_timeout = 0;

--bun:split
DROP TABLE IF EXISTS seller_earnings;

--bun:split
DROP TABLE IF EXISTS bid_transactions;

--bun:split
DROP TABLE IF EXISTS auction_bid_holds;

--bun:split
DROP TABLE IF EXISTS auction_bids;

--bun:split
ALTER TABLE auctions
DROP COLUMN IF EXISTS winner_id,
DROP COLUMN IF EXISTS settled_at;

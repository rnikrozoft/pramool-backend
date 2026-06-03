SET statement_timeout = 0;

--bun:split

ALTER TABLE bid_transactions
    DROP CONSTRAINT IF EXISTS bid_transactions_auction_id_fkey;

--bun:split

DELETE FROM bid_transactions WHERE auction_id IS NULL;

--bun:split

ALTER TABLE bid_transactions
    ALTER COLUMN auction_id SET NOT NULL;

--bun:split

ALTER TABLE bid_transactions
    ADD CONSTRAINT bid_transactions_auction_id_fkey
    FOREIGN KEY (auction_id) REFERENCES auctions(auction_id) ON DELETE CASCADE;

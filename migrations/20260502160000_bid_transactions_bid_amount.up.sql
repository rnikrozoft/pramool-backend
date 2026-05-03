-- Optional bid price / refund display amount for wallet activity UI.
ALTER TABLE bid_transactions ADD COLUMN IF NOT EXISTS bid_amount BIGINT;

SET statement_timeout = 0;

--bun:split

DROP INDEX IF EXISTS idx_transactions_dispute_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS credit_reversed;
ALTER TABLE transactions DROP COLUMN IF EXISTS dispute_status;
ALTER TABLE transactions DROP COLUMN IF EXISTS dispute_id;

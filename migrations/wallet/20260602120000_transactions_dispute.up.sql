SET statement_timeout = 0;

--bun:split

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS dispute_id VARCHAR(255);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS dispute_status VARCHAR(30);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS credit_reversed BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_transactions_dispute_id ON transactions (dispute_id) WHERE dispute_id IS NOT NULL;

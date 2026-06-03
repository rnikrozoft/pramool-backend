SET statement_timeout = 0;

--bun:split

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS suspend_fulfillment_pending BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_users_suspend_fulfillment_pending
    ON users (suspend_fulfillment_pending)
    WHERE suspend_fulfillment_pending = TRUE;

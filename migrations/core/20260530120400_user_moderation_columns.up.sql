SET statement_timeout = 0;

--bun:split

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS restricted_until TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS restricted_reason TEXT,
    ADD COLUMN IF NOT EXISTS buyer_reputation_points BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS suspended_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS suspended_reason TEXT,
    ADD COLUMN IF NOT EXISTS suspended_by_admin_id INT REFERENCES admin_users (admin_id);

CREATE INDEX IF NOT EXISTS idx_users_restricted_until
    ON users (restricted_until)
    WHERE restricted_until IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_suspended_at
    ON users (suspended_at)
    WHERE suspended_at IS NOT NULL;

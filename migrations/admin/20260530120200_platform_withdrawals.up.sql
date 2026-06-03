SET statement_timeout = 0;

--bun:split

CREATE TABLE platform_withdrawals (
    withdrawal_id BIGSERIAL PRIMARY KEY,
    admin_id INT NOT NULL REFERENCES admin_users (admin_id),
    amount_baht BIGINT NOT NULL,
    omise_transfer_id TEXT,
    omise_status TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_platform_withdrawals_created ON platform_withdrawals (created_at DESC);

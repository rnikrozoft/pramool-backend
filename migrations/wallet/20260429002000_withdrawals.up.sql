SET statement_timeout = 0;

--bun:split

CREATE TABLE withdrawals (
    withdrawal_id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    fee_amount BIGINT NOT NULL DEFAULT 0,
    transfer_amount BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    omise_transfer_id TEXT,
    omise_recipient_id TEXT,
    balance_before BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--bun:split

CREATE INDEX idx_withdrawals_user_created_at ON withdrawals(user_id, created_at DESC);

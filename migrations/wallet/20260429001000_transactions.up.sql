SET statement_timeout = 0;

--bun:split

CREATE TABLE transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    charge_id VARCHAR(255) NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    fee_amount BIGINT NOT NULL DEFAULT 0,
    credit_amount BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    paid BOOLEAN NOT NULL DEFAULT FALSE,
    credited BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

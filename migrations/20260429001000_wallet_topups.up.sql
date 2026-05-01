SET statement_timeout = 0;

--bun:split
CREATE TABLE wallet_topups (
    charge_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(13) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    paid BOOLEAN NOT NULL DEFAULT FALSE,
    credited BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

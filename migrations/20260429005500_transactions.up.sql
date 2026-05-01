SET statement_timeout = 0;

--bun:split
CREATE TABLE IF NOT EXISTS transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    charge_id VARCHAR(255) NOT NULL UNIQUE,
    user_id VARCHAR(13) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    paid BOOLEAN NOT NULL DEFAULT FALSE,
    credited BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--bun:split
INSERT INTO transactions (charge_id, user_id, amount, status, paid, credited, created_at, updated_at)
SELECT charge_id, user_id, amount, status, paid, credited, created_at, updated_at
FROM wallet_topups
ON CONFLICT (charge_id) DO NOTHING;

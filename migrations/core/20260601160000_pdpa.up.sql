SET statement_timeout = 0;

--bun:split

ALTER TABLE tel_verify
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

--bun:split

CREATE TABLE consent_log (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users (user_id) ON DELETE SET NULL,
    tel VARCHAR(10),
    consent_type VARCHAR(50) NOT NULL,
    policy_version VARCHAR(32) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_consent_log_user_id ON consent_log (user_id);
CREATE INDEX idx_consent_log_tel ON consent_log (tel);
CREATE INDEX idx_consent_log_created_at ON consent_log (created_at);

--bun:split

CREATE TABLE dsar_requests (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    request_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    user_note TEXT,
    admin_note TEXT,
    handled_by_admin_id INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_dsar_requests_user_id ON dsar_requests (user_id);
CREATE INDEX idx_dsar_requests_status ON dsar_requests (status);
CREATE INDEX idx_dsar_requests_created_at ON dsar_requests (created_at);

SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS dsar_request_exports (
    dsar_id BIGINT PRIMARY KEY REFERENCES dsar_requests (id) ON DELETE CASCADE,
    export_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

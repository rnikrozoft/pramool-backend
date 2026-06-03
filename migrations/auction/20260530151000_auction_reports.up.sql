SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS auction_reports (
    report_id BIGSERIAL PRIMARY KEY,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions (auction_id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    reporter_user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    admin_id INT REFERENCES admin_users (admin_id),
    admin_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    CONSTRAINT auction_reports_status_check CHECK (status IN ('pending', 'accepted', 'rejected'))
);

CREATE INDEX IF NOT EXISTS idx_auction_reports_status_created ON auction_reports (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auction_reports_auction ON auction_reports (auction_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_auction_reports_pending_reporter
    ON auction_reports (auction_id, reporter_user_id)
    WHERE status = 'pending';

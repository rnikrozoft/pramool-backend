SET statement_timeout = 0;

--bun:split

CREATE TABLE user_restriction_appeals (
    appeal_id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    admin_id INT REFERENCES admin_users (admin_id),
    admin_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    CONSTRAINT user_restriction_appeals_status_check CHECK (status IN ('pending', 'accepted', 'rejected'))
);

CREATE INDEX idx_user_restriction_appeals_status_created ON user_restriction_appeals (status, created_at DESC);

CREATE UNIQUE INDEX idx_user_restriction_appeals_pending_user
    ON user_restriction_appeals (user_id)
    WHERE status = 'pending';

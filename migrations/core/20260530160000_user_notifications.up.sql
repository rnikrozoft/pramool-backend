SET statement_timeout = 0;

--bun:split

CREATE TABLE user_notifications (
    notification_id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    kind VARCHAR(40) NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    read_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_notifications_user_expires ON user_notifications (user_id, expires_at DESC);
CREATE INDEX idx_user_notifications_user_unread ON user_notifications (user_id) WHERE read_at IS NULL;

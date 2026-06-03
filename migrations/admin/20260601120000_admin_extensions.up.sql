SET statement_timeout = 0;

--bun:split

ALTER TABLE auctions
    ADD COLUMN IF NOT EXISTS admin_ship_deadline_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS admin_confirm_deadline_at TIMESTAMPTZ;

--bun:split

CREATE TABLE IF NOT EXISTS platform_settings (
    setting_key   TEXT PRIMARY KEY,
    setting_value JSONB NOT NULL DEFAULT '{}',
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by_admin_id INT REFERENCES admin_users(admin_id)
);

--bun:split

INSERT INTO platform_settings (setting_key, setting_value) VALUES
    ('fulfillment', '{"seller_ship_deadline_days":7,"escrow_auto_confirm_days":7}'::jsonb),
    ('fees', '{"seller_keep_normal_pct":90,"seller_keep_early_pct":95,"min_withdraw_credit_thb":100,"omise_transfer_fee_thb":30}'::jsonb),
    ('notifications', '{"slack_webhook_url":"","notify_email":"","last_notify_at":null}'::jsonb)
ON CONFLICT (setting_key) DO NOTHING;

--bun:split

CREATE TABLE IF NOT EXISTS site_announcements (
    announcement_id BIGSERIAL PRIMARY KEY,
    title           TEXT NOT NULL,
    body            TEXT NOT NULL DEFAULT '',
    link_url        TEXT,
    severity        TEXT NOT NULL DEFAULT 'info' CHECK (severity IN ('info', 'warning', 'critical')),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    starts_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ends_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by_admin_id INT REFERENCES admin_users(admin_id)
);

--bun:split

CREATE INDEX IF NOT EXISTS idx_site_announcements_active ON site_announcements (is_active, starts_at DESC);

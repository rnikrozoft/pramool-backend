SET statement_timeout = 0;

--bun:split

INSERT INTO platform_settings (setting_key, setting_value) VALUES
    ('listing', '{
        "min_start_price_thb": 100,
        "bid_cancel_option_fee_thb": 1,
        "auto_renew_option_fee_thb": 20,
        "free_listing_duration_days": 2,
        "extra_listing_day_fee_pct": 1,
        "auto_renew_success_fee_pct": 1
    }'::jsonb)
ON CONFLICT (setting_key) DO NOTHING;

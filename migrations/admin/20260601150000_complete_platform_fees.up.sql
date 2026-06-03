SET statement_timeout = 0;

--bun:split

UPDATE platform_settings
SET setting_value = setting_value || '{
    "seller_keep_normal_pct": 75,
    "seller_keep_early_pct": 70,
    "min_withdraw_credit_thb": 100,
    "omise_transfer_fee_thb": 21,
    "min_topup_gross_thb": 100,
    "omise_promptpay_fee_ppm": 17655
}'::jsonb
WHERE setting_key = 'fees';

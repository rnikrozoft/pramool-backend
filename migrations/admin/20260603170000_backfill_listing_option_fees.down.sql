SET statement_timeout = 0;

--bun:split

UPDATE platform_settings
SET setting_value = setting_value - 'auto_renew_option_fee_thb' - 'bid_cancel_option_fee_thb'
WHERE setting_key = 'listing';

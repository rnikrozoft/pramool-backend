SET statement_timeout = 0;

--bun:split

UPDATE platform_settings
SET setting_value = setting_value || '{"auto_renew_option_fee_thb": 20}'::jsonb
WHERE setting_key = 'listing';

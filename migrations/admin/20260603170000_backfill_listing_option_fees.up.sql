SET statement_timeout = 0;

--bun:split

-- Backfill listing option fees when row predates field or admin save omitted keys.
UPDATE platform_settings
SET setting_value = setting_value || '{"auto_renew_option_fee_thb": 20}'::jsonb
WHERE setting_key = 'listing'
  AND NOT (setting_value ? 'auto_renew_option_fee_thb');

--bun:split

UPDATE platform_settings
SET setting_value = setting_value || '{"bid_cancel_option_fee_thb": 1}'::jsonb
WHERE setting_key = 'listing'
  AND NOT (setting_value ? 'bid_cancel_option_fee_thb');

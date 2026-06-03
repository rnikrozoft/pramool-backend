SET statement_timeout = 0;

--bun:split

UPDATE platform_settings
SET setting_value = setting_value
    || '{"seller_keep_normal_pct":75,"seller_keep_early_pct":70}'::jsonb
WHERE setting_key = 'fees'
  AND COALESCE(setting_value->>'seller_keep_normal_pct', '90') = '90';

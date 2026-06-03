SET statement_timeout = 0;

--bun:split

DROP TABLE IF EXISTS site_announcements;

--bun:split

DROP TABLE IF EXISTS platform_settings;

--bun:split

ALTER TABLE auctions
    DROP COLUMN IF EXISTS admin_ship_deadline_at,
    DROP COLUMN IF EXISTS admin_confirm_deadline_at;

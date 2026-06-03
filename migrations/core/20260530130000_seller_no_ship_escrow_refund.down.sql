SET statement_timeout = 0;

--bun:split

ALTER TABLE users DROP COLUMN IF EXISTS seller_no_ship_count;

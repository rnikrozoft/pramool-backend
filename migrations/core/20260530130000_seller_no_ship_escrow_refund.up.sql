SET statement_timeout = 0;

--bun:split

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS seller_no_ship_count INT NOT NULL DEFAULT 0;

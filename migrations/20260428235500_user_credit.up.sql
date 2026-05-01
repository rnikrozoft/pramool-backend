SET statement_timeout = 0;

--bun:split
ALTER TABLE users
ADD COLUMN credit BIGINT NOT NULL DEFAULT 0;

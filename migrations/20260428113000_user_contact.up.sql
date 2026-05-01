SET statement_timeout = 0;

--bun:split
ALTER TABLE users
ADD COLUMN email VARCHAR(120) DEFAULT '',
ADD COLUMN facebook VARCHAR(255) DEFAULT '';

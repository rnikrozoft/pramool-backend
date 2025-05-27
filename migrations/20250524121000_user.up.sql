SET statement_timeout = 0;

--bun:split

CREATE TABLE users (
    user_id VARCHAR(13) PRIMARY KEY,
    tel VARCHAR(10) NOT NULL,
    first_name VARCHAR(20) NOT NULL,
    last_name VARCHAR(20) NOT NULL,
    address_primary TEXT NOT NULL,
    address TEXT,
    soi TEXT,
    road TEXT,
    sub_district TEXT NOT NULL,
    district TEXT NOT NULL,
    province TEXT NOT NULL,
    zip_code VARCHAR(5) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


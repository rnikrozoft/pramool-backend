SET statement_timeout = 0;

--bun:split

CREATE TABLE users (
    user_id VARCHAR(13) PRIMARY KEY,
    tel VARCHAR(10) NOT NULL REFERENCES tel_verify(tel),
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
    email VARCHAR(120) NOT NULL DEFAULT '',
    facebook VARCHAR(255) NOT NULL DEFAULT '',
    credit BIGINT NOT NULL DEFAULT 0,
    password_hash VARCHAR(255),
    bank_id BIGINT REFERENCES banks(bank_id),
    bank_account_name TEXT,
    bank_account_number VARCHAR(20),
    omise_recipient_id TEXT,
    seller_review_points_total BIGINT NOT NULL DEFAULT 0,
    seller_review_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

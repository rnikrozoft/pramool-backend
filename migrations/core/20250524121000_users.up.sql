SET statement_timeout = 0;

--bun:split

CREATE EXTENSION IF NOT EXISTS pgcrypto;

--bun:split

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_id_hash VARCHAR(64) UNIQUE,
    national_id_enc TEXT,
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
    deleted_at TIMESTAMPTZ,
    anonymized_at TIMESTAMPTZ,
    marketing_opt_in BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_users_national_id_hash ON users (national_id_hash) WHERE national_id_hash IS NOT NULL;
CREATE INDEX idx_users_deleted_at ON users (deleted_at) WHERE deleted_at IS NOT NULL;

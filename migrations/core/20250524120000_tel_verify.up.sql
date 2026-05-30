SET statement_timeout = 0;

--bun:split

CREATE TABLE tel_verify (
    tel VARCHAR(10) PRIMARY KEY,
    otp_timeout_count INTEGER NOT NULL DEFAULT 0,
    otp_banned_until TIMESTAMPTZ,
    signup_first_name VARCHAR(100),
    signup_last_name VARCHAR(100),
    signup_email VARCHAR(120) NOT NULL DEFAULT '',
    password_hash VARCHAR(255)
);

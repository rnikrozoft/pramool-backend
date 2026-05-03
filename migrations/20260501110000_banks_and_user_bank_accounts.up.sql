SET statement_timeout = 0;

--bun:split
CREATE TABLE IF NOT EXISTS banks (
    bank_id BIGSERIAL PRIMARY KEY,
    bank_code VARCHAR(10) NOT NULL UNIQUE,
    bank_name_th VARCHAR(100) NOT NULL,
    bank_name_en VARCHAR(100) NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--bun:split
INSERT INTO banks (bank_code, bank_name_th, bank_name_en, display_order, is_active)
VALUES
    ('KBANK', 'กสิกรไทย', 'Kasikornbank', 1, TRUE),
    ('KTB', 'กรุงไทย', 'Krung Thai Bank', 2, TRUE),
    ('BBL', 'กรุงเทพ', 'Bangkok Bank', 3, TRUE),
    ('SCB', 'ไทยพาณิชย์', 'Siam Commercial Bank', 4, TRUE),
    ('BAY', 'กรุงศรีอยุธยา', 'Bank of Ayudhya', 5, TRUE),
    ('TTB', 'ทีทีบี', 'ttb', 6, TRUE),
    ('GSB', 'ออมสิน', 'Government Savings Bank', 7, TRUE),
    ('BAAC', 'ธ.ก.ส.', 'Bank for Agriculture and Agricultural Cooperatives', 8, TRUE),
    ('CIMBT', 'ซีไอเอ็มบี ไทย', 'CIMB Thai Bank', 9, TRUE),
    ('UOB', 'ยูโอบี', 'United Overseas Bank', 10, TRUE)
ON CONFLICT (bank_code) DO UPDATE SET
    bank_name_th = EXCLUDED.bank_name_th,
    bank_name_en = EXCLUDED.bank_name_en,
    display_order = EXCLUDED.display_order,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

--bun:split
ALTER TABLE users
ADD COLUMN IF NOT EXISTS bank_id BIGINT REFERENCES banks(bank_id),
ADD COLUMN IF NOT EXISTS bank_account_name TEXT,
ADD COLUMN IF NOT EXISTS bank_account_number VARCHAR(20);

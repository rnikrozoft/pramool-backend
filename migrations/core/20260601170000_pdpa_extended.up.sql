SET statement_timeout = 0;

--bun:split

ALTER TABLE dsar_requests
    ADD COLUMN IF NOT EXISTS due_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS deletion_executed_at TIMESTAMPTZ;

UPDATE dsar_requests SET due_at = created_at + INTERVAL '30 days' WHERE due_at IS NULL;

--bun:split

CREATE TABLE data_disclosure_log (
    id BIGSERIAL PRIMARY KEY,
    disclosure_type VARCHAR(50) NOT NULL,
    actor_user_id UUID NOT NULL,
    subject_user_id UUID NOT NULL,
    reference_id VARCHAR(64),
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_data_disclosure_log_actor ON data_disclosure_log (actor_user_id);
CREATE INDEX idx_data_disclosure_log_subject ON data_disclosure_log (subject_user_id);
CREATE INDEX idx_data_disclosure_log_created_at ON data_disclosure_log (created_at);

--bun:split

CREATE TABLE data_processors (
    processor_id SERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    purpose TEXT NOT NULL,
    data_categories TEXT NOT NULL,
    location VARCHAR(120) NOT NULL DEFAULT 'Thailand',
    privacy_url TEXT,
    dpa_status VARCHAR(40) NOT NULL DEFAULT 'pending',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO data_processors (name, purpose, data_categories, location, privacy_url, dpa_status, sort_order) VALUES
('Omise Co., Ltd.', 'ชำระเงิน PromptPay เติมเครดิต และโอนเงินถอนเครดิต', 'ชื่อบัญชี เลขบัญชี จำนวนเงิน รหัสอ้างอิงธุรกรรม', 'Thailand', 'https://www.omise.co/privacy', 'active', 1),
('Thai Bulk SMS / ผู้ให้บริการ SMS OTP', 'ส่ง OTP ยืนยันเบอร์โทรศัพท์', 'เบอร์โทรศัพท์ ข้อความ OTP', 'Thailand', NULL, 'active', 2),
('TrackingMore', 'ติดตามสถานะพัสดุจากเลข tracking', 'เลขพัสดุ รหัสขนส่ง สถานะการจัดส่ง', 'Singapore / EU (cloud)', 'https://www.trackingmore.com/privacy-policy', 'active', 3);

--bun:split

CREATE TABLE privacy_breach_incidents (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    severity VARCHAR(20) NOT NULL DEFAULT 'medium',
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    discovered_at TIMESTAMPTZ NOT NULL,
    reported_to_dpc_at TIMESTAMPTZ,
    affected_user_estimate INT,
    summary TEXT NOT NULL,
    remediation TEXT,
    created_by_admin_id INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ
);

CREATE INDEX idx_privacy_breach_incidents_status ON privacy_breach_incidents (status);

--bun:split

CREATE TABLE dpia_records (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    processing_activity TEXT NOT NULL,
    risk_level VARCHAR(20) NOT NULL DEFAULT 'medium',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    review_notes TEXT,
    next_review_at DATE,
    created_by_admin_id INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--bun:split

CREATE TABLE retention_job_definitions (
    job_id VARCHAR(64) PRIMARY KEY,
    name_th VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    retention_days INT NOT NULL,
    batch_runner VARCHAR(40) NOT NULL DEFAULT 'future',
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    last_run_at TIMESTAMPTZ,
    last_deleted_count BIGINT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO retention_job_definitions (job_id, name_th, description, retention_days, batch_runner, is_enabled) VALUES
('tel_verify_stale', 'ลบ tel_verify ค้าง', 'ลบข้อมูลสมัครค้างที่ยังไม่สมัครครบและไม่มี user', 30, 'inline', TRUE),
('user_notifications_expired', 'ลบการแจ้งเตือนหมดอายุ', 'ลบ user_notifications ที่ expires_at เก่ากว่ากำหนด', 90, 'future', TRUE),
('consent_log_archive', 'จัดเก็บ consent log เก่า', 'ย้าย/ลบ consent_log เก่ากว่านโยบาย (รอ batch runner)', 2555, 'future', FALSE),
('access_log_files', 'ลบ access log ไฟล์', 'ลบ structured access log เก่ากว่า 12 เดือน (รอ batch runner / log infra)', 365, 'future', FALSE),
('data_disclosure_log', 'ลบ disclosure log เก่า', 'ลบ data_disclosure_log เก่ากว่า 24 เดือน (รอ batch runner)', 730, 'future', FALSE);

SET statement_timeout = 0;

--bun:split

CREATE TABLE admin_sensitive_access_log (
    id BIGSERIAL PRIMARY KEY,
    admin_id INT NOT NULL,
    user_id UUID NOT NULL,
    field_type VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admin_sensitive_access_log_admin_id ON admin_sensitive_access_log (admin_id);
CREATE INDEX idx_admin_sensitive_access_log_user_id ON admin_sensitive_access_log (user_id);
CREATE INDEX idx_admin_sensitive_access_log_created_at ON admin_sensitive_access_log (created_at);

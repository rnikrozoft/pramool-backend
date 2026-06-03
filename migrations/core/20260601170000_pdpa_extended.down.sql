DELETE FROM retention_job_definitions;
DROP TABLE IF EXISTS retention_job_definitions;
DROP TABLE IF EXISTS dpia_records;
DROP TABLE IF EXISTS privacy_breach_incidents;
DROP TABLE IF EXISTS data_processors;
DROP TABLE IF EXISTS data_disclosure_log;

ALTER TABLE dsar_requests
    DROP COLUMN IF EXISTS deletion_executed_at,
    DROP COLUMN IF EXISTS due_at;

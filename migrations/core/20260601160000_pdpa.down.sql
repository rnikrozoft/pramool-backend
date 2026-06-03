DROP TABLE IF EXISTS dsar_requests;
DROP TABLE IF EXISTS consent_log;

ALTER TABLE tel_verify
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS created_at;

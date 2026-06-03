DROP INDEX IF EXISTS idx_users_posting_restricted_until;

ALTER TABLE users
    DROP COLUMN IF EXISTS posting_restricted_until,
    DROP COLUMN IF EXISTS posting_restricted_reason;

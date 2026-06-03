DROP INDEX IF EXISTS idx_users_suspended_at;
DROP INDEX IF EXISTS idx_users_restricted_until;

ALTER TABLE users
    DROP COLUMN IF EXISTS suspended_by_admin_id,
    DROP COLUMN IF EXISTS suspended_reason,
    DROP COLUMN IF EXISTS suspended_at,
    DROP COLUMN IF EXISTS buyer_reputation_points,
    DROP COLUMN IF EXISTS restricted_reason,
    DROP COLUMN IF EXISTS restricted_until;

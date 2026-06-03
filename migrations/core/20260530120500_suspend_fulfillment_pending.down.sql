DROP INDEX IF EXISTS idx_users_suspend_fulfillment_pending;

ALTER TABLE users DROP COLUMN IF EXISTS suspend_fulfillment_pending;

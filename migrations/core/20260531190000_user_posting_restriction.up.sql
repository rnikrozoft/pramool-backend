SET statement_timeout = 0;

--bun:split

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS posting_restricted_until TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS posting_restricted_reason TEXT;

--bun:split

CREATE INDEX IF NOT EXISTS idx_users_posting_restricted_until
    ON users (posting_restricted_until)
    WHERE posting_restricted_until IS NOT NULL;

--bun:split

-- Report-accept used full account restriction; move sellers to posting-only and clear bidders.
UPDATE users u
SET
    posting_restricted_until = u.restricted_until,
    posting_restricted_reason = u.restricted_reason,
    restricted_until = NULL,
    restricted_reason = NULL,
    updated_at = NOW()
FROM (
    SELECT DISTINCT seller_id AS user_id
    FROM auction_reports
    WHERE status = 'accepted'
) s
WHERE u.user_id = s.user_id
  AND u.restricted_until IS NOT NULL
  AND u.restricted_until > NOW()
  AND COALESCE(u.restricted_reason, '') LIKE '%รายงานผู้ใช้%';

--bun:split

UPDATE users
SET restricted_until = NULL,
    restricted_reason = NULL,
    updated_at = NOW()
WHERE restricted_until IS NOT NULL
  AND restricted_until > NOW()
  AND COALESCE(restricted_reason, '') LIKE '%รายงานผู้ใช้%';

SET statement_timeout = 0;

--bun:split

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
  AND u.restricted_until > NOW();

--bun:split

UPDATE users
SET restricted_until = NULL,
    restricted_reason = NULL,
    updated_at = NOW()
WHERE restricted_until IS NOT NULL
  AND restricted_until > NOW()
  AND COALESCE(restricted_reason, '') LIKE '%ฝ่าฝืนกฎ%';

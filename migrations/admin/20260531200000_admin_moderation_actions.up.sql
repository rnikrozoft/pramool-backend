SET statement_timeout = 0;

--bun:split

CREATE TABLE admin_moderation_actions (
    action_id BIGSERIAL PRIMARY KEY,
    admin_id INT NOT NULL REFERENCES admin_users(admin_id),
    action_type VARCHAR(40) NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    target_id VARCHAR(40) NOT NULL,
    summary TEXT NOT NULL,
    note TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--bun:split

CREATE INDEX idx_admin_moderation_actions_created ON admin_moderation_actions(created_at DESC);

--bun:split

CREATE INDEX idx_admin_moderation_actions_type_created ON admin_moderation_actions(action_type, created_at DESC);

--bun:split

INSERT INTO admin_moderation_actions (admin_id, action_type, target_type, target_id, summary, note, metadata, created_at)
SELECT
    ar.admin_id,
    CASE WHEN ar.status = 'accepted' THEN 'report_accept' ELSE 'report_reject' END,
    'auction',
    ar.auction_id,
    CASE
        WHEN ar.status = 'accepted'
            THEN 'รับร้องเรียน ลบรายการ ' || ar.auction_id || ' (ข้อมูลย้อนหลัง)'
        ELSE 'ปฏิเสธร้องเรียน รายการ ' || ar.auction_id || ' (ข้อมูลย้อนหลัง)'
    END,
    MAX(ar.admin_note),
    jsonb_build_object('auction_id', ar.auction_id, 'seller_id', (array_agg(ar.seller_id))[1], 'backfill', true),
    ar.resolved_at
FROM auction_reports ar
WHERE ar.status IN ('accepted', 'rejected')
  AND ar.admin_id IS NOT NULL
  AND ar.resolved_at IS NOT NULL
GROUP BY ar.admin_id, ar.status, ar.auction_id, ar.resolved_at;

--bun:split

INSERT INTO admin_moderation_actions (admin_id, action_type, target_type, target_id, summary, note, metadata, created_at)
SELECT
    a.admin_id,
    CASE WHEN a.status = 'accepted' THEN 'appeal_accept' ELSE 'appeal_reject' END,
    'appeal',
    a.appeal_id::text,
    CASE
        WHEN a.status = 'accepted'
            THEN 'อนุมัติอุทธรณ์ ผู้ใช้ ' || a.user_id || ' (ข้อมูลย้อนหลัง)'
        ELSE 'ปฏิเสธอุทธรณ์ ผู้ใช้ ' || a.user_id || ' (ข้อมูลย้อนหลัง)'
    END,
    a.admin_note,
    jsonb_build_object('appeal_id', a.appeal_id, 'user_id', a.user_id, 'backfill', true),
    a.resolved_at
FROM user_restriction_appeals a
WHERE a.status IN ('accepted', 'rejected')
  AND a.admin_id IS NOT NULL
  AND a.resolved_at IS NOT NULL;

--bun:split

INSERT INTO admin_moderation_actions (admin_id, action_type, target_type, target_id, summary, note, metadata, created_at)
SELECT
    u.suspended_by_admin_id,
    'user_suspend',
    'user',
    u.user_id,
    'ระงับบัญชี ' || u.user_id || ' (ข้อมูลย้อนหลัง)',
    u.suspended_reason,
    jsonb_build_object('user_id', u.user_id, 'backfill', true),
    u.suspended_at
FROM users u
WHERE u.suspended_at IS NOT NULL
  AND u.suspended_by_admin_id IS NOT NULL;

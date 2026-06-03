SET statement_timeout = 0;

--bun:split

-- Backfill in-app notifications for report acceptances that ran before
-- notifyViolatingAuctionRemovalTx was deployed in backoffice-core.

INSERT INTO user_notifications (user_id, kind, title, body, expires_at)
SELECT DISTINCT ar.seller_id,
    'listing_removed_seller',
    'รายการประมูลถูกลบ',
    'รายการประมูลของคุณถูกลบ บัญชีถูกจำกัดการใช้งาน 30 วัน (เติม/ถอน/โพส/บิด) และคะแนนรีวิวผู้ขาย −10 คะแนน' || E'\n' ||
    'เหตุผล: รายการประมูลฝ่าฝืนกฎ (รายงานผู้ใช้)',
    TIMESTAMPTZ '2099-12-31 23:59:59+00'
FROM auction_reports ar
WHERE ar.status = 'accepted'
  AND NOT EXISTS (
      SELECT 1
      FROM user_notifications n
      WHERE n.user_id = ar.seller_id
        AND n.kind = 'listing_removed_seller'
  );

--bun:split

INSERT INTO user_notifications (user_id, kind, title, body, expires_at)
SELECT DISTINCT ar.reporter_user_id,
    'report_reward',
    'รางวัลการร้องเรียน',
    'คำร้องเรียนของคุณได้รับการอนุมัติ ได้รับเครดิต 1 ฿ เข้ากระเป๋าแล้ว',
    TIMESTAMPTZ '2099-12-31 23:59:59+00'
FROM auction_reports ar
WHERE ar.status = 'accepted'
  AND NOT EXISTS (
      SELECT 1
      FROM user_notifications n
      WHERE n.user_id = ar.reporter_user_id
        AND n.kind = 'report_reward'
  );

--bun:split

INSERT INTO user_notifications (user_id, kind, title, body, expires_at)
SELECT u.user_id,
    'listing_removed_bidder',
    'การจำกัดบัญชีจากรายการฝ่าฝืน',
    'มีรายการประมูลที่คุณเข้าร่วมถูกลบ เงินประมูลถูกคืนแล้ว บัญชีถูกจำกัด 30 วัน และคะแนนผู้ประมูล −5 คะแนน' || E'\n' ||
    'เหตุผล: ' || COALESCE(u.restricted_reason, 'รายการประมูลฝ่าฝืนกฎ (รายงานผู้ใช้)'),
    TIMESTAMPTZ '2099-12-31 23:59:59+00'
FROM users u
WHERE u.restricted_until > NOW()
  AND COALESCE(u.restricted_reason, '') LIKE '%รายงานผู้ใช้%'
  AND NOT EXISTS (
      SELECT 1
      FROM user_notifications n
      WHERE n.user_id = u.user_id
        AND n.kind IN ('listing_removed_seller', 'listing_removed_bidder')
  );

SET statement_timeout = 0;

--bun:split
ALTER TABLE auctions
    ADD COLUMN IF NOT EXISTS seller_shipped_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS buyer_received_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS seller_payout_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS payout_early_close BOOLEAN NOT NULL DEFAULT FALSE;

--bun:split
-- ประมูลที่ปิดไปแล้วและมีผู้ชนะ: ถือว่าจ่ายผู้ขายครบแล้ว (ก่อนมี escrow)
UPDATE auctions
SET seller_payout_at = settled_at
WHERE status = 'closed'
  AND COALESCE(NULLIF(TRIM(winner_id), ''), '') <> ''
  AND settled_at IS NOT NULL
  AND seller_payout_at IS NULL;

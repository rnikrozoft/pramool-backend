SET statement_timeout = 0;

-- รองรับหลายหมวดหมู่คั่นด้วย | (สูงสุด 5 รายการจากแอป)
ALTER TABLE auctions
    ALTER COLUMN category TYPE TEXT;

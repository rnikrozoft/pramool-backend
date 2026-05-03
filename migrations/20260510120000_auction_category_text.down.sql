-- คืนค่า VARCHAR — แถวที่ยาวเกิน 100 ตัวอักษรจะล้มเหลว ควรย่อข้อความก่อนถอย migration
ALTER TABLE auctions
    ALTER COLUMN category TYPE VARCHAR(100);

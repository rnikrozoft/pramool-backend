SET statement_timeout = 0;

--bun:split

ALTER TABLE product_categories
    ADD COLUMN IF NOT EXISTS is_system BOOLEAN NOT NULL DEFAULT FALSE;

--bun:split

UPDATE product_categories
SET is_system = TRUE
WHERE name = 'อื่นๆ';

--bun:split

INSERT INTO product_categories (name, sort_order, is_active, is_system)
SELECT 'อื่นๆ', 80, TRUE, TRUE
WHERE NOT EXISTS (SELECT 1 FROM product_categories WHERE name = 'อื่นๆ');

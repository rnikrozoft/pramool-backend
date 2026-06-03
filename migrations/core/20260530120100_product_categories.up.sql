SET statement_timeout = 0;

--bun:split

CREATE TABLE product_categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_categories_active_sort ON product_categories (is_active, sort_order, name);

--bun:split

INSERT INTO product_categories (name, sort_order) VALUES
    ('เครื่องใช้ไฟฟ้า', 10),
    ('โทรศัพท์มือถือ', 20),
    ('แท็บเล็ต', 30),
    ('คอมพิวเตอร์', 40),
    ('กล้องถ่ายรูป', 50),
    ('แฟชั่น', 60),
    ('ของสะสม', 70),
    ('อื่นๆ', 80),
    ('เกมคอนโซล', 90),
    ('กระเป๋า', 100);

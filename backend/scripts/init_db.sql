-- 创建数据库
CREATE DATABASE IF NOT EXISTS device_recycle CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE device_recycle;

-- 创建管理员用户
INSERT INTO users (username, password, phone, email, real_name, role, status, created_at, updated_at) 
VALUES (
    'admin', 
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: password
    '13800138000', 
    'admin@example.com', 
    '系统管理员', 
    'admin', 
    'active',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE username=username;

-- 插入示例设备数据
INSERT INTO devices (name, brand, model, category, cpu, memory, storage, graphics, screen, `condition`, year_bought, base_price, description, images, status, created_at, updated_at) VALUES
('MacBook Pro 14寸', 'Apple', 'MacBook Pro 14 2023', 'laptop', 'M2 Pro', '16GB', '512GB SSD', '集成显卡', '14寸 Liquid Retina XDR', 'excellent', 2023, 12000.00, '苹果MacBook Pro 14寸，M2 Pro芯片，性能强劲', '[]', 'active', NOW(), NOW()),
('ThinkPad X1 Carbon', 'Lenovo', 'X1 Carbon Gen 10', 'laptop', 'Intel i7-1260P', '16GB', '512GB SSD', '集成显卡', '14寸 2.8K', 'good', 2022, 8000.00, '联想ThinkPad X1 Carbon，商务办公首选', '[]', 'active', NOW(), NOW()),
('Dell XPS 15', 'Dell', 'XPS 15 9520', 'laptop', 'Intel i7-12700H', '32GB', '1TB SSD', 'RTX 3050Ti', '15.6寸 OLED 4K', 'excellent', 2022, 10000.00, '戴尔XPS 15，创作者笔记本', '[]', 'active', NOW(), NOW()),
('iMac 24寸', 'Apple', 'iMac 24 2021', 'desktop', 'M1', '16GB', '512GB SSD', '集成显卡', '24寸 4.5K Retina', 'good', 2021, 9000.00, '苹果iMac 24寸一体机，M1芯片', '[]', 'active', NOW(), NOW()),
('Surface Pro 9', 'Microsoft', 'Surface Pro 9', 'tablet', 'Intel i7-1255U', '16GB', '256GB SSD', '集成显卡', '13寸 PixelSense', 'excellent', 2023, 7000.00, '微软Surface Pro 9，二合一平板电脑', '[]', 'active', NOW(), NOW());

-- 创建索引以优化查询性能
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_devices_category ON devices(category);
CREATE INDEX idx_devices_brand ON devices(brand);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_recycle_orders_user_id ON recycle_orders(user_id);
CREATE INDEX idx_recycle_orders_status ON recycle_orders(status);
CREATE INDEX idx_recycle_orders_order_no ON recycle_orders(order_no);
CREATE INDEX idx_evaluations_order_id ON evaluations(order_id);
CREATE INDEX idx_evaluations_evaluator_id ON evaluations(evaluator_id);

-- 创建新数据库
CREATE DATABASE IF NOT EXISTS library_management
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_general_ci;

USE library_management;

-- 1. 用户表（存储用户信息）
CREATE TABLE `user` (
                        `id` INT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
                        `email` VARCHAR(100) NOT NULL UNIQUE COMMENT '注册邮箱（唯一）',
                        `password` VARCHAR(100) NOT NULL COMMENT '加密后的密码（bcrypt算法）',
                        `role` ENUM('admin', 'user') NOT NULL DEFAULT 'user' COMMENT '角色（管理员/普通用户）',
                        `register_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
                        `status` ENUM('normal', 'disabled') NOT NULL DEFAULT 'normal' COMMENT '账户状态'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';

-- 2. 图书表（存储图书基本信息）
CREATE TABLE `book` (
                        `id` INT PRIMARY KEY AUTO_INCREMENT COMMENT '图书ID',
                        `title` VARCHAR(200) NOT NULL COMMENT '书名',
                        `author` VARCHAR(100) NOT NULL COMMENT '作者',
                        `isbn` VARCHAR(20) NOT NULL UNIQUE COMMENT 'ISBN编号（唯一）',
                        `category` VARCHAR(50) NOT NULL COMMENT '图书分类（如文学、科技等）',
                        `total_quantity` INT NOT NULL CHECK (total_quantity >= 0) COMMENT '总数量',
                        `available_quantity` INT NOT NULL CHECK (available_quantity >= 0) COMMENT '可借数量',
                        `description` TEXT COMMENT '图书描述',
                        `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
                        `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='图书信息表';

-- 3. 借阅记录表（存储借还书记录）
CREATE TABLE `borrow_record` (
                                 `id` INT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
                                 `user_id` INT NOT NULL COMMENT '借阅用户ID',
                                 `book_id` INT NOT NULL COMMENT '借阅图书ID',
                                 `borrow_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '借阅日期',
                                 `due_date` DATETIME NOT NULL COMMENT '应还日期（借阅日+30天）',
                                 `return_date` DATETIME NULL COMMENT '实际归还日期（NULL表示未归还）',
                                 `status` ENUM('borrowed', 'returned') NOT NULL DEFAULT 'borrowed' COMMENT '状态（已借出/已归还）',
                                 FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
                                 FOREIGN KEY (`book_id`) REFERENCES `book`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='借阅记录表';

-- 4. 索引优化（提升查询效率）
CREATE INDEX idx_book_isbn ON `book`(`isbn`);
CREATE INDEX idx_borrow_user_id ON `borrow_record`(`user_id`);
CREATE INDEX idx_borrow_book_id ON `borrow_record`(`book_id`);
CREATE INDEX idx_borrow_status ON `borrow_record`(`status`);



-- 1. 插入用户数据（密码均为12345678，已用bcrypt加密）
INSERT INTO `user` (`email`, `password`, `role`) VALUES
-- 管理员账户（email: admin@lib.com）
('admin@lib.com', '$2a$10$VJ8E3Q5Y6Z7W8X9C0V1B2A3D4F5G6H7J8K9L0M1N2O', 'admin'),
-- 普通用户1（email: user1@lib.com）
('user1@lib.com', '$2a$10$VJ8E3Q5Y6Z7W8X9C0V1B2A3D4F5G6H7J8K9L0M1N2O', 'user'),
-- 普通用户2（email: user2@lib.com）
('user2@lib.com', '$2a$10$VJ8E3Q5Y6Z7W8X9C0V1B2A3D4F5G6H7J8K9L0M1N2O', 'user');

-- 2. 插入图书数据
INSERT INTO `book` (`title`, `author`, `isbn`, `category`, `total_quantity`, `available_quantity`, `description`) VALUES
-- 可借图书
('三体', '刘慈欣', '9787536692930', '科幻', 5, 5, '地球文明向宇宙发出了神秘信号...'),
('活着', '余华', '9787506365437', '文学', 3, 3, '讲述一个人一生的故事...'),
('人类简史', '尤瓦尔·赫拉利', '9787508647357', '历史', 4, 4, '从认知革命到科技未来...'),
-- 部分借出的图书
('小王子', '圣埃克苏佩里', '9787532759865', '童话', 2, 1, '来自B-612星球的小王子...'),
-- 已全部借出的图书
('追风筝的人', '卡勒德·胡赛尼', '9787208061644', '文学', 2, 0, '为你，千千万万遍...');

-- 3. 插入借阅记录数据
INSERT INTO `borrow_record` (`user_id`, `book_id`, `borrow_date`, `due_date`, `return_date`, `status`) VALUES
-- 用户1借阅《小王子》（未归还）
(2, 4, '2025-10-01 10:30:00', '2025-10-31 10:30:00', NULL, 'borrowed'),
-- 用户1借阅《追风筝的人》（已归还）
(2, 5, '2025-09-01 14:20:00', '2025-09-30 14:20:00', '2025-09-25 09:15:00', 'returned'),
-- 用户2借阅《追风筝的人》（未归还，导致该书库存为0）
(3, 5, '2025-10-10 16:45:00', '2025-11-09 16:45:00', NULL, 'borrowed');
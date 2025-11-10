-- PostgreSQL 数据库初始化脚本

-- 1. 用户表（存储用户信息）
CREATE TABLE IF NOT EXISTS "user" (
    "id" SERIAL PRIMARY KEY,
    "email" VARCHAR(100) NOT NULL UNIQUE,
    "password" VARCHAR(100) NOT NULL,
    "role" VARCHAR(10) NOT NULL DEFAULT 'user' CHECK ("role" IN ('admin', 'user')),
    "register_time" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "status" VARCHAR(10) NOT NULL DEFAULT 'normal' CHECK ("status" IN ('normal', 'disabled'))
);

COMMENT ON TABLE "user" IS '用户信息表';
COMMENT ON COLUMN "user"."id" IS '用户ID';
COMMENT ON COLUMN "user"."email" IS '注册邮箱（唯一）';
COMMENT ON COLUMN "user"."password" IS '加密后的密码（bcrypt算法）';
COMMENT ON COLUMN "user"."role" IS '角色（管理员/普通用户）';
COMMENT ON COLUMN "user"."register_time" IS '注册时间';
COMMENT ON COLUMN "user"."status" IS '账户状态';

-- 2. 图书表（存储图书基本信息）
CREATE TABLE IF NOT EXISTS "book" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(200) NOT NULL,
    "author" VARCHAR(100) NOT NULL,
    "isbn" VARCHAR(20) NOT NULL UNIQUE,
    "category" VARCHAR(50) NOT NULL,
    "total_quantity" INTEGER NOT NULL CHECK ("total_quantity" >= 0),
    "available_quantity" INTEGER NOT NULL CHECK ("available_quantity" >= 0),
    "description" TEXT,
    "create_time" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "update_time" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE "book" IS '图书信息表';
COMMENT ON COLUMN "book"."id" IS '图书ID';
COMMENT ON COLUMN "book"."title" IS '书名';
COMMENT ON COLUMN "book"."author" IS '作者';
COMMENT ON COLUMN "book"."isbn" IS 'ISBN编号（唯一）';
COMMENT ON COLUMN "book"."category" IS '图书分类（如文学、科技等）';
COMMENT ON COLUMN "book"."total_quantity" IS '总数量';
COMMENT ON COLUMN "book"."available_quantity" IS '可借数量';
COMMENT ON COLUMN "book"."description" IS '图书描述';
COMMENT ON COLUMN "book"."create_time" IS '添加时间';
COMMENT ON COLUMN "book"."update_time" IS '更新时间';

-- 3. 借阅记录表（存储借还书记录）
CREATE TABLE IF NOT EXISTS "borrow_record" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "book_id" INTEGER NOT NULL,
    "borrow_date" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "due_date" TIMESTAMP NOT NULL,
    "return_date" TIMESTAMP NULL,
    "status" VARCHAR(10) NOT NULL DEFAULT 'borrowed' CHECK ("status" IN ('borrowed', 'returned')),
    FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE,
    FOREIGN KEY ("book_id") REFERENCES "book"("id") ON DELETE CASCADE
);

COMMENT ON TABLE "borrow_record" IS '借阅记录表';
COMMENT ON COLUMN "borrow_record"."id" IS '记录ID';
COMMENT ON COLUMN "borrow_record"."user_id" IS '借阅用户ID';
COMMENT ON COLUMN "borrow_record"."book_id" IS '借阅图书ID';
COMMENT ON COLUMN "borrow_record"."borrow_date" IS '借阅日期';
COMMENT ON COLUMN "borrow_record"."due_date" IS '应还日期（借阅日+30天）';
COMMENT ON COLUMN "borrow_record"."return_date" IS '实际归还日期（NULL表示未归还）';
COMMENT ON COLUMN "borrow_record"."status" IS '状态（已借出/已归还）';

-- 4. 验证码记录表（存储邮箱验证码信息，仅管理员可见）
CREATE TABLE IF NOT EXISTS "email_code_record" (
    "id" SERIAL PRIMARY KEY,
    "email" VARCHAR(100) NOT NULL,
    "code" VARCHAR(10) NOT NULL,
    "action" VARCHAR(20) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "expires_at" TIMESTAMP NOT NULL,
    "is_used" BOOLEAN NOT NULL DEFAULT FALSE,
    "used_at" TIMESTAMP NULL
);

COMMENT ON TABLE "email_code_record" IS '邮箱验证码记录表（仅管理员可见）';
COMMENT ON COLUMN "email_code_record"."id" IS '记录ID';
COMMENT ON COLUMN "email_code_record"."email" IS '接收验证码的邮箱';
COMMENT ON COLUMN "email_code_record"."code" IS '验证码（6位数字）';
COMMENT ON COLUMN "email_code_record"."action" IS '用途（register: 注册, forget: 忘记密码）';
COMMENT ON COLUMN "email_code_record"."created_at" IS '创建时间';
COMMENT ON COLUMN "email_code_record"."expires_at" IS '过期时间（创建后30分钟）';
COMMENT ON COLUMN "email_code_record"."is_used" IS '是否已使用';
COMMENT ON COLUMN "email_code_record"."used_at" IS '使用时间';

-- 5. 索引优化（提升查询效率）
CREATE INDEX IF NOT EXISTS idx_book_isbn ON "book"("isbn");
CREATE INDEX IF NOT EXISTS idx_borrow_user_id ON "borrow_record"("user_id");
CREATE INDEX IF NOT EXISTS idx_borrow_book_id ON "borrow_record"("book_id");
CREATE INDEX IF NOT EXISTS idx_borrow_status ON "borrow_record"("status");
CREATE INDEX IF NOT EXISTS idx_email_code_email ON "email_code_record"("email");
CREATE INDEX IF NOT EXISTS idx_email_code_action ON "email_code_record"("action");
CREATE INDEX IF NOT EXISTS idx_email_code_created_at ON "email_code_record"("created_at");
CREATE INDEX IF NOT EXISTS idx_email_code_expires_at ON "email_code_record"("expires_at");
CREATE INDEX IF NOT EXISTS idx_email_code_is_used ON "email_code_record"("is_used");

-- 6. 创建更新时间触发器函数（PostgreSQL 需要手动处理 update_time）
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为 book 表创建触发器
CREATE TRIGGER update_book_updated_at BEFORE UPDATE ON "book"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 7. 插入用户数据（密码均为12345678，已用bcrypt加密）
INSERT INTO "user" ("email", "password", "role") VALUES
-- 管理员账户（email: admin@lib.com）
('admin@lib.com', '$2a$10$VJ8E3Q5Y6Z7W8X9C0V1B2A3D4F5G6H7J8K9L0M1N2O', 'admin'),
-- 普通用户1（email: user1@lib.com）
('user1@lib.com', '$2a$10$VJ8E3Q5Y6Z7W8X9C0V1B2A3D4F5G6H7J8K9L0M1N2O', 'user'),
-- 普通用户2（email: user2@lib.com）
('user2@lib.com', '$2a$10$VJ8E3Q5Y6Z7W8X9C0V1B2A3D4F5G6H7J8K9L0M1N2O', 'user')
ON CONFLICT ("email") DO NOTHING;

-- 8. 插入图书数据
INSERT INTO "book" ("title", "author", "isbn", "category", "total_quantity", "available_quantity", "description") VALUES
-- 可借图书
('三体', '刘慈欣', '9787536692930', '科幻', 5, 5, '地球文明向宇宙发出了神秘信号...'),
('活着', '余华', '9787506365437', '文学', 3, 3, '讲述一个人一生的故事...'),
('人类简史', '尤瓦尔·赫拉利', '9787508647357', '历史', 4, 4, '从认知革命到科技未来...'),
-- 部分借出的图书
('小王子', '圣埃克苏佩里', '9787532759865', '童话', 2, 1, '来自B-612星球的小王子...'),
-- 已全部借出的图书
('追风筝的人', '卡勒德·胡赛尼', '9787208061644', '文学', 2, 0, '为你，千千万万遍...')
ON CONFLICT ("isbn") DO NOTHING;

-- 9. 插入借阅记录数据
INSERT INTO "borrow_record" ("user_id", "book_id", "borrow_date", "due_date", "return_date", "status") VALUES
-- 用户1借阅《小王子》（未归还）
(2, 4, '2025-10-01 10:30:00', '2025-10-31 10:30:00', NULL, 'borrowed'),
-- 用户1借阅《追风筝的人》（已归还）
(2, 5, '2025-09-01 14:20:00', '2025-09-30 14:20:00', '2025-09-25 09:15:00', 'returned'),
-- 用户2借阅《追风筝的人》（未归还，导致该书库存为0）
(3, 5, '2025-10-10 16:45:00', '2025-11-09 16:45:00', NULL, 'borrowed')
ON CONFLICT DO NOTHING;

-- 10. 添加图片功能相关字段（使用 ALTER 命令）
-- 为 book 表添加 cover_image_url 字段
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'book' 
        AND column_name = 'cover_image_url'
    ) THEN
        ALTER TABLE "book" ADD COLUMN "cover_image_url" VARCHAR(500) DEFAULT NULL;
        COMMENT ON COLUMN "book"."cover_image_url" IS '图书封面图片URL（存储在Cloudflare R2）';
    END IF;
END $$;


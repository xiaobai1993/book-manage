-- 添加验证码记录表到现有数据库
-- 适用于本地Docker PostgreSQL环境

-- 验证码记录表（存储邮箱验证码信息，仅管理员可见）
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

-- 创建索引（提升查询效率）
CREATE INDEX IF NOT EXISTS idx_email_code_email ON "email_code_record"("email");
CREATE INDEX IF NOT EXISTS idx_email_code_action ON "email_code_record"("action");
CREATE INDEX IF NOT EXISTS idx_email_code_created_at ON "email_code_record"("created_at");
CREATE INDEX IF NOT EXISTS idx_email_code_expires_at ON "email_code_record"("expires_at");
CREATE INDEX IF NOT EXISTS idx_email_code_is_used ON "email_code_record"("is_used");

-- 验证表是否创建成功
SELECT 'email_code_record表创建成功！' AS message;


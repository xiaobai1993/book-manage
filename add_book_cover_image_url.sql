-- 为 book 表添加 cover_image_url 字段
-- PostgreSQL 版本

-- 如果字段已存在，不会报错（使用 IF NOT EXISTS）
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


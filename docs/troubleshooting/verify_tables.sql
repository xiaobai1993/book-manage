-- 检查 Supabase 数据库中表是否存在
-- 在 Supabase SQL Editor 中执行此脚本

-- 1. 查看所有表
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_type = 'BASE TABLE'
ORDER BY table_name;

-- 2. 检查 user 表是否存在
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_schema = 'public' 
    AND table_name = 'user'
) AS user_table_exists;

-- 3. 检查 book 表是否存在
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_schema = 'public' 
    AND table_name = 'book'
) AS book_table_exists;

-- 4. 检查 borrow_record 表是否存在
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_schema = 'public' 
    AND table_name = 'borrow_record'
) AS borrow_record_table_exists;

-- 5. 如果表存在，查看 user 表的结构
SELECT column_name, data_type, is_nullable
FROM information_schema.columns
WHERE table_schema = 'public' 
AND table_name = 'user'
ORDER BY ordinal_position;

-- 6. 尝试查询 user 表（如果存在）
SELECT COUNT(*) as user_count FROM "user";




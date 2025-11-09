# 修复数据库表不存在错误

## 🔴 错误信息

```
ERROR: relation "user" does not exist (SQLSTATE 42P01)
```

## 📝 问题原因

这个错误表示 Supabase 数据库中还没有创建表。需要在 Supabase 中执行 SQL 脚本来创建表结构。

## ✅ 解决方案

### 步骤 1：登录 Supabase Dashboard

1. 访问 [Supabase Dashboard](https://app.supabase.com)
2. 选择你的项目

### 步骤 2：打开 SQL Editor

1. 在左侧菜单中，点击 **SQL Editor**
2. 点击 **New query** 按钮

### 步骤 3：执行 SQL 脚本

1. 打开项目中的 `data_postgresql.sql` 文件
2. **复制全部内容**（从第 1 行到最后一行）
3. 粘贴到 Supabase SQL Editor 中
4. 点击 **Run** 按钮（或按 `Ctrl+Enter` / `Cmd+Enter`）

### 步骤 4：验证表创建成功

执行完成后，你应该看到：

1. **成功消息**：显示执行成功
2. **查看表**：
   - 在左侧菜单点击 **Table Editor**
   - 应该看到三个表：
     - `user` - 用户表
     - `book` - 图书表
     - `borrow_record` - 借阅记录表

### 步骤 5：验证数据

在 **Table Editor** 中：

1. 点击 `user` 表
2. 应该看到一些初始数据（如 `admin@lib.com` 等）
3. 点击 `book` 表
4. 应该看到一些示例图书数据

## 🔍 如果执行失败

### 问题 1：表已存在错误

如果看到 "relation already exists" 错误：

1. **删除现有表**（如果表是空的）：
   ```sql
   DROP TABLE IF EXISTS "borrow_record";
   DROP TABLE IF EXISTS "book";
   DROP TABLE IF EXISTS "user";
   ```

2. 然后重新执行 `data_postgresql.sql` 脚本

### 问题 2：触发器错误

如果触发器创建失败，可以忽略（不影响基本功能），或者单独执行：

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_book_updated_at BEFORE UPDATE ON "book"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 问题 3：数据插入失败

如果数据插入失败（如密码哈希值错误），可以：

1. **跳过数据插入部分**：只执行表创建部分（前 88 行）
2. **手动插入管理员账户**：
   ```sql
   INSERT INTO "user" ("email", "password", "role") VALUES
   ('admin@lib.com', '$2a$10$VJ8E3Q5Y6Z7W8X9C0V1B2A3D4F5G6H7J8K9L0M1N2O', 'admin')
   ON CONFLICT ("email") DO NOTHING;
   ```

## ✅ 验证修复

修复后，重新测试注册功能：

1. 在 Render 日志中，应该不再看到 "relation does not exist" 错误
2. 注册接口应该可以正常工作
3. 可以成功创建新用户

## 📋 快速检查清单

- [ ] 已登录 Supabase Dashboard
- [ ] 已打开 SQL Editor
- [ ] 已复制 `data_postgresql.sql` 全部内容
- [ ] 已执行 SQL 脚本
- [ ] 在 Table Editor 中看到三个表
- [ ] 表中有初始数据
- [ ] Render 日志中不再有错误
- [ ] 注册功能正常工作

## 🔗 相关文档

- [部署指南](./DEPLOY.md) - 查看完整的部署步骤
- [Supabase 连接修复](./SUPABASE_CONNECTION_FIX.md) - 数据库连接问题




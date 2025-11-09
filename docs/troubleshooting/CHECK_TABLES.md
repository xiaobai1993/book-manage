# 检查 Supabase 数据库表是否存在

## 🔍 快速检查步骤

### 方法一：通过 Supabase Dashboard 检查

1. **登录 Supabase Dashboard**
   - 访问 https://app.supabase.com
   - 选择你的项目

2. **查看 Table Editor**
   - 左侧菜单 → **Table Editor**
   - 查看是否有以下三个表：
     - `user`
     - `book`
     - `borrow_record`

3. **如果看不到表**，说明表确实不存在，需要执行 SQL 脚本创建

### 方法二：通过 SQL Editor 查询

1. **打开 SQL Editor**
   - 左侧菜单 → **SQL Editor**
   - 点击 **New query**

2. **执行查询语句**：
   ```sql
   SELECT table_name 
   FROM information_schema.tables 
   WHERE table_schema = 'public' 
   AND table_type = 'BASE TABLE'
   ORDER BY table_name;
   ```

3. **查看结果**
   - 如果结果为空或只有系统表，说明用户表不存在
   - 如果看到 `user`、`book`、`borrow_record`，说明表存在

### 方法三：直接查询 user 表

在 SQL Editor 中执行：

```sql
SELECT * FROM "user" LIMIT 1;
```

- **如果成功**：说明表存在，可能是其他问题（如权限、schema）
- **如果报错 "relation does not exist"**：说明表确实不存在

## 🔧 如果表确实不存在

如果通过以上方法确认表不存在，需要执行 `data_postgresql.sql` 脚本：

1. 在 SQL Editor 中
2. 复制 `data_postgresql.sql` 的全部内容
3. 粘贴并执行

## 🔍 如果表存在但报错

如果表确实存在但应用报错，可能是：

1. **Schema 问题**：表在错误的 schema 中
2. **权限问题**：连接用户没有访问权限
3. **连接字符串问题**：连接到了错误的数据库

检查连接字符串中的数据库名：
- 你的连接字符串：`postgresql://...@...:6543/postgres`
- 数据库名是 `postgres`，这是默认数据库

如果表在其他数据库中，需要修改连接字符串。




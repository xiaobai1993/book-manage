# Supabase 连接问题解决方案

## 🔴 问题：network is unreachable

如果你看到类似错误：
```
dial tcp [2406:da18:...]:5432: connect: network is unreachable
```

这通常是因为：
1. Supabase 需要使用 **Connection Pooling** 连接字符串
2. 或者需要配置 IP 白名单

## ✅ 解决方案

### 方法一：使用 Connection Pooling（推荐）

Supabase 提供了两种连接方式：
1. **Direct connection** - 直接连接（可能被限制）
2. **Connection Pooling** - 连接池（推荐用于生产环境）

#### 步骤：

1. **在 Supabase Dashboard 中**：
   - Settings → Database
   - 找到 **Connection Pooling** 部分
   - 选择 **Session** 模式（或 **Transaction** 模式）
   - 复制 **Connection string**（URI 格式）

2. **Connection Pooling 的连接字符串格式**：
   ```
   postgresql://postgres.xxx:DbPw87Jk2xRn93Qs@aws-0-xxx.pooler.supabase.com:6543/postgres
   ```
   
   注意：
   - 主机名是 `xxx.pooler.supabase.com`（不是 `db.xxx.supabase.co`）
   - 端口是 `6543`（不是 `5432`）
   - 用户名是 `postgres.xxx`（包含项目引用）

3. **在 Render 中设置**：
   - 使用 Connection Pooling 的连接字符串作为 `DATABASE_URL`

### 方法二：配置 IP 白名单

如果必须使用直接连接：

1. **在 Supabase Dashboard 中**：
   - Settings → Database
   - 找到 **Connection Pooling** 或 **Network Restrictions**
   - 添加 Render 的 IP 地址（或允许所有 IP）

2. **Render 的 IP 地址**：
   - Render 使用动态 IP，建议允许所有 IP 或使用 Connection Pooling

### 方法三：使用 Connection Pooling 的 Transaction 模式

如果 Session 模式不行，尝试 Transaction 模式：

1. 在 Supabase → Settings → Database → Connection Pooling
2. 选择 **Transaction** 模式
3. 复制连接字符串
4. 在 Render 中使用

## 📝 连接字符串对比

### Direct Connection（直接连接）
```
postgresql://postgres:密码@db.xxx.supabase.co:5432/postgres
```
- 端口：5432
- 可能被网络限制

### Connection Pooling（连接池）
```
postgresql://postgres.xxx:密码@aws-0-xxx.pooler.supabase.com:6543/postgres
```
- 端口：6543（Session）或 5432（Transaction）
- 推荐用于生产环境
- 更好的性能和稳定性

## 🔧 快速修复步骤

1. **获取 Connection Pooling 连接字符串**：
   - Supabase Dashboard → Settings → Database
   - Connection Pooling → Session mode
   - 复制 URI 连接字符串

2. **更新 Render 环境变量**：
   - 将 `DATABASE_URL` 替换为 Connection Pooling 的连接字符串

3. **重新部署**：
   - Render 会自动重新部署
   - 查看日志确认连接成功

## ✅ 验证

部署成功后，日志应该显示：
```
Database connection established successfully (type: postgres)
```

如果没有这个信息，检查：
- 连接字符串是否正确
- 密码是否正确（特殊字符需要 URL 编码）
- Supabase 项目是否正常运行


# 注册接口超时问题排查指南

## 🔍 问题现象

点击注册时，接口提示超时（Timeout）。

## 📋 排查步骤

### 1. 检查 Render 后端日志

1. 登录 [Render Dashboard](https://dashboard.render.com)
2. 进入你的后端服务
3. 点击 **Logs** 标签页
4. 查看最近的错误日志

**重点关注**：
- 是否有数据库连接错误
- 是否有超时错误
- 查看注册接口的详细日志（已添加性能日志）

### 2. 检查数据库连接

#### 2.1 验证 Supabase 连接字符串

确保使用的是 **Connection Pooling** 连接字符串：

1. 登录 [Supabase Dashboard](https://app.supabase.com)
2. Settings → Database → Connection Pooling
3. 选择 **Session** 模式
4. 复制连接字符串

**连接字符串格式**：
```
postgresql://postgres.xxx:密码@xxx.pooler.supabase.com:6543/postgres?sslmode=require
```

**注意**：
- 端口应该是 `6543`（Session 模式）或 `5432`（Transaction 模式）
- 主机名应该是 `xxx.pooler.supabase.com`（不是 `db.xxx.supabase.co`）
- 必须包含 `sslmode=require`

#### 2.2 检查 Render 环境变量

在 Render Dashboard 中，检查以下环境变量：

- `DATABASE_URL` - 应该是完整的连接字符串
- 数据库类型 - 仅支持 PostgreSQL（DB_TYPE 环境变量已不再使用）

### 3. 测试数据库连接

在 Render 日志中，应该看到：
```
Database connection established successfully (type: postgres)
Connection pool: MaxOpen=25, MaxIdle=10
```

如果没有看到这些日志，说明数据库连接失败。

### 4. 检查注册接口性能日志

注册接口已添加详细的时间日志，在 Render 日志中应该看到：

```
[Register] 邮箱检查耗时: XXXms
[Register] 密码加密耗时: XXXms
[Register] 创建用户耗时: XXXms
[Register] 总耗时: XXXms
```

**正常耗时参考**：
- 邮箱检查：< 100ms
- 密码加密：< 500ms
- 创建用户：< 200ms
- 总耗时：< 1000ms

如果某个步骤耗时过长，说明问题出在该步骤。

### 5. 常见问题及解决方案

#### 问题 1：数据库连接超时

**症状**：
- 日志显示 "failed to connect database" 或 "timeout"
- 注册接口直接超时

**解决方案**：
1. 确认使用 Connection Pooling 连接字符串
2. 检查 Supabase 项目是否正常运行
3. 检查网络连接（Render 到 Supabase 的网络）

#### 问题 2：密码加密耗时过长

**症状**：
- 日志显示 "密码加密耗时" > 2秒

**解决方案**：
- bcrypt 加密本身比较耗时，这是正常的
- 如果超过 5 秒，可能是服务器资源不足
- 考虑升级 Render 服务计划

#### 问题 3：数据库查询超时

**症状**：
- 日志显示 "邮箱检查耗时" 或 "创建用户耗时" > 5秒

**解决方案**：
1. 检查 Supabase 数据库性能
2. 确认数据库表结构正确（已创建索引）
3. 检查是否有数据库连接池耗尽

#### 问题 4：Render 服务休眠

**症状**：
- 首次请求超时，后续请求正常
- 日志显示服务正在启动

**解决方案**：
- Render 免费计划会在 15 分钟无活动后休眠
- 首次请求需要约 30 秒唤醒时间
- 考虑升级到付费计划（$7/月）避免休眠

### 6. 快速诊断命令

如果可以通过 Render 的 Shell 访问，可以运行：

```bash
# 测试数据库连接
psql $DATABASE_URL -c "SELECT 1;"
```

### 7. 临时解决方案

如果问题紧急，可以：

1. **增加超时时间**（前端）：
   - 检查前端 API 请求的超时设置
   - 临时增加超时时间到 30 秒

2. **添加重试机制**（前端）：
   - 如果超时，自动重试一次

3. **优化数据库连接**：
   - 已添加连接池配置
   - 确保使用 Connection Pooling

### 8. 验证修复

修复后，测试注册流程：

1. 发送验证码
2. 输入验证码和用户信息
3. 点击注册
4. 检查是否成功
5. 查看 Render 日志确认没有错误

## 📝 日志示例

### 正常日志
```
Database connection established successfully (type: postgres)
Connection pool: MaxOpen=25, MaxIdle=10
[Register] 邮箱检查耗时: 45ms
[Register] 密码加密耗时: 234ms
[Register] 创建用户耗时: 67ms
[Register] 总耗时: 346ms
```

### 异常日志
```
[Register] 邮箱检查耗时: 5000ms  ← 数据库查询慢
[Register] 密码加密耗时: 234ms
[Register] 创建用户失败: timeout, 耗时: 10000ms  ← 数据库写入超时
```

## 🔗 相关文档

- [部署指南](./DEPLOY.md)
- [Supabase 连接修复](./SUPABASE_CONNECTION_FIX.md)
- [Render 文档](https://render.com/docs)




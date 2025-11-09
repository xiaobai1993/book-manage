# QQ 邮箱发送验证码超时问题排查

## 🔴 问题现象

发送邮件验证码时接口超时，使用 QQ 邮箱。

## 📋 排查步骤

### 1. 检查 Render 日志

1. 登录 [Render Dashboard](https://dashboard.render.com)
2. 进入后端服务 → **Logs**
3. 查看邮件发送相关的日志

**应该看到的日志**：
```
[Email Service] 开始发送验证码到 xxx@example.com (action: register)
[Email Service] 开始发送邮件 (SSL, 465端口) 到 xxx@example.com
[Email Service] SSL连接成功 (耗时: XXXms)
[Email Service] SMTP认证成功 (耗时: XXXms)
[Email Service] 邮件内容发送成功 (耗时: XXXms)
[Email Service] 邮件发送总耗时: XXXms
[Email Service] 成功发送验证码邮件到 xxx@example.com (register)
```

**如果看到错误**：
- `连接SMTP服务器失败` - 网络连接问题
- `SMTP认证失败` - 用户名或密码错误
- `timeout` - 连接超时

### 2. 检查 QQ 邮箱配置

#### 2.1 确认使用授权码

**重要**：QQ 邮箱必须使用**授权码**，不能使用登录密码！

1. 登录 [QQ 邮箱](https://mail.qq.com)
2. 点击 **设置** → **账户**
3. 找到 **POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务**
4. 开启 **POP3/SMTP服务** 或 **IMAP/SMTP服务**
5. 点击 **生成授权码**
6. **复制授权码**（16位字符）

#### 2.2 检查 Render 环境变量

在 Render Dashboard 中，检查以下环境变量：

| 变量名 | 值 | 说明 |
|--------|-----|------|
| `SMTP_HOST` | `smtp.qq.com` | QQ 邮箱 SMTP 服务器 |
| `SMTP_PORT` | `465` | SSL 端口（推荐）或 `587`（STARTTLS） |
| `SMTP_USER` | `你的QQ号@qq.com` | 完整的 QQ 邮箱地址 |
| `SMTP_PASSWORD` | `授权码` | **必须是授权码，不是登录密码** |

**常见错误**：
- ❌ `SMTP_PASSWORD` 使用登录密码
- ❌ `SMTP_USER` 只填写 QQ 号（缺少 @qq.com）
- ❌ `SMTP_PORT` 使用错误的端口

### 3. 检查网络连接

Render 服务器可能无法访问 QQ 邮箱的 SMTP 服务器（网络限制）。

**测试方法**：
1. 查看 Render 日志中的连接错误
2. 如果看到 `connection refused` 或 `timeout`，可能是网络问题

**解决方案**：
- 使用其他邮件服务（如 Gmail、SendGrid、Mailgun）
- 或使用邮件服务代理

### 4. 检查端口配置

QQ 邮箱支持两种方式：

#### 方式一：SSL（465端口）- 推荐
```
SMTP_HOST=smtp.qq.com
SMTP_PORT=465
```
- 使用 SSL 加密连接
- 代码中已实现 `sendEmailSSL` 方法

#### 方式二：STARTTLS（587端口）
```
SMTP_HOST=smtp.qq.com
SMTP_PORT=587
```
- 使用 STARTTLS 加密
- 需要确保代码支持（当前代码支持）

### 5. 查看详细日志

代码已添加详细的性能日志，在 Render 日志中应该看到：

```
[Email Service] 开始发送验证码到 xxx@example.com (action: register)
[Email Service] 开始发送邮件 (SSL, 465端口) 到 xxx@example.com
[Email Service] SSL连接成功 (耗时: 234ms)
[Email Service] SMTP认证成功 (耗时: 123ms)
[Email Service] 邮件内容发送成功 (耗时: 456ms)
[Email Service] 连接关闭成功 (耗时: 12ms)
[Email Service] 邮件发送总耗时: 825ms
[Email Service] 成功发送验证码邮件到 xxx@example.com (register)
[SendEmailCode] 发送验证码成功 (耗时: 830ms)
```

**正常耗时参考**：
- SSL 连接：< 1秒
- SMTP 认证：< 1秒
- 发送内容：< 2秒
- 总耗时：< 5秒

如果某个步骤耗时过长，说明问题出在该步骤。

## 🔧 常见问题及解决方案

### 问题 1：SMTP 认证失败

**错误信息**：
```
[Email Service] SMTP认证失败: 535 Error: authentication failed
```

**原因**：
- 使用了登录密码而不是授权码
- 授权码错误
- 用户名格式错误

**解决方案**：
1. 确认使用授权码（不是登录密码）
2. 重新生成授权码
3. 确认 `SMTP_USER` 是完整的邮箱地址（如 `123456789@qq.com`）

### 问题 2：连接超时

**错误信息**：
```
[Email Service] SSL连接失败: timeout
```

**原因**：
- Render 服务器无法访问 QQ 邮箱 SMTP 服务器
- 网络限制或防火墙

**解决方案**：
1. **使用其他邮件服务**（推荐）：
   - Gmail（需要应用专用密码）
   - SendGrid（免费额度：100封/天）
   - Mailgun（免费额度：5000封/月）
   - 阿里云邮件推送

2. **使用邮件服务代理**：
   - 通过代理服务器转发邮件

3. **临时方案**：
   - 在开发环境测试，生产环境使用其他服务

### 问题 3：端口连接被拒绝

**错误信息**：
```
[Email Service] SSL连接失败: connection refused
```

**原因**：
- 端口配置错误
- QQ 邮箱服务未开启

**解决方案**：
1. 确认 `SMTP_PORT` 是 `465` 或 `587`
2. 确认 QQ 邮箱已开启 SMTP 服务
3. 尝试切换端口（465 ↔ 587）

### 问题 4：发送成功但收不到邮件

**可能原因**：
1. 邮件被放入垃圾箱
2. 收件邮箱地址错误
3. QQ 邮箱限制（发送频率过高）

**解决方案**：
1. 检查垃圾箱
2. 确认收件邮箱地址正确
3. 降低发送频率（代码中已限制 1 分钟内不能重复发送）

## 🚀 快速修复步骤

### 步骤 1：确认授权码

1. 登录 QQ 邮箱
2. 设置 → 账户 → 生成授权码
3. 复制授权码

### 步骤 2：更新 Render 环境变量

在 Render Dashboard 中更新：

```
SMTP_HOST=smtp.qq.com
SMTP_PORT=465
SMTP_USER=你的QQ号@qq.com
SMTP_PASSWORD=你的授权码（16位）
```

### 步骤 3：重新部署

1. Render 会自动重新部署
2. 或手动触发部署

### 步骤 4：测试

1. 在前端点击"发送验证码"
2. 查看 Render 日志
3. 检查邮箱（包括垃圾箱）

## 📝 使用其他邮件服务

如果 QQ 邮箱无法使用，可以切换到其他服务：

### Gmail

```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=465
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=应用专用密码（不是登录密码）
```

### SendGrid（推荐，免费）

1. 注册 [SendGrid](https://sendgrid.com)
2. 创建 API Key
3. 使用 SMTP 方式：

```
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=你的API Key
```

## ✅ 验证修复

修复后，测试发送验证码：

1. 前端点击"发送验证码"
2. 查看 Render 日志，应该看到成功日志
3. 检查邮箱（包括垃圾箱）
4. 验证码应该能正常收到

## 🔗 相关文档

- [部署指南](./DEPLOY.md)
- [超时问题排查](./TROUBLESHOOTING_TIMEOUT.md)



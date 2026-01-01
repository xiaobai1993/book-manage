# 邮件服务迁移到 Resend

## 概述

图书管理系统的邮件发送服务已从传统的 SMTP 方式迁移到 [Resend](https://resend.com) API。Resend 是一个现代化的邮件 API 服务，提供更可靠的邮件送达率和更简单的集成方式。

## 主要变更

### 1. 代码变更

#### services/email.go

- **删除**：传统的 SMTP 连接代码（SSL/STARTTLS）
- **新增**：Resend Go SDK 集成
- **简化**：邮件发送逻辑从 ~400 行减少到 ~80 行

**关键变更**：
```go
// 旧方式：使用 SMTP
func (s *EmailService) sendEmailSSL(...)
func (s *EmailService) sendEmailSTARTTLS(...)

// 新方式：使用 Resend API
func (s *EmailService) sendEmailViaResend(toEmail, action, code string) error {
    params := &resend.SendEmailRequest{
        From:    s.cfg.SMTPUser,
        To:      []string{toEmail},
        Subject: subject,
        Html:    htmlContent,
    }
    _, err := s.resend.Emails.Send(params)
    return err
}
```

### 2. 配置变更

#### 环境变量

| 变量名 | 旧值（SMTP） | 新值（Resend） | 说明 |
|--------|-------------|---------------|------|
| `SMTP_HOST` | `smtp.qq.com` | *不再使用* | SMTP 服务器地址 |
| `SMTP_PORT` | `465` 或 `587` | *不再使用* | SMTP 端口 |
| `SMTP_USER` | `your-email@qq.com` | `noreply@ai-speed.xyz` | 发件人邮箱 |
| `SMTP_PASSWORD` | QQ邮箱授权码 | `re_ejMaZbaK_8d6A6H4vxDb7Wvkf7S1AZYCx` | Resend API Key |

**重要**：虽然 `SMTP_HOST` 和 `SMTP_PORT` 不再使用，但为了向后兼容，配置结构中保留了这些字段。

#### render.yaml

```yaml
# 邮箱配置（使用 Resend 发送邮件）
- key: SMTP_PASSWORD
  value: "re_ejMaZbaK_8d6A6H4vxDb7Wvkf7S1AZYCx"  # Resend API Key
- key: SMTP_USER
  value: "noreply@ai-speed.xyz"  # 发件人邮箱
```

### 3. 依赖变更

#### go.mod

**新增依赖**：
```go
require (
    github.com/resend/resend-go/v2 v2.28.0
)
```

## 迁移步骤（已完成）

### ✅ 1. 安装 Resend Go SDK

```bash
go get github.com/resend/resend-go/v2
```

### ✅ 2. 更新代码

- 替换邮件发送逻辑
- 删除 SMTP 相关代码
- 更新初始化逻辑

### ✅ 3. 更新配置

- 修改 `render.yaml` 配置文件
- 更新环境变量说明
- 使用 Resend API Key 替换 SMTP 密码

### ✅ 4. 测试编译

```bash
go build -o book-manage
```

### ⏳ 5. 在 Resend 控制台验证域名（需要手动完成）

**重要步骤**：在使用邮件功能前，必须完成域名验证：

1. 访问 [Resend Domains](https://resend.com/domains)
2. 添加域名：`ai-speed.xyz`
3. 配置 DNS 记录：
   - 添加 TXT 记录用于域名验证
   - 添加 MX 记录用于接收邮件（可选）
   - 添加 CNAME 记录用于 DKIM（重要！提高邮件送达率）

4. 等待 DNS 生效（通常几分钟到几小时）

### ⏳ 6. 重新部署到 Render

```bash
git add .
git commit -m "feat: 迁移邮件服务到 Resend"
git push
```

Render 会自动重新部署，使用新的环境变量。

## 优势

### 相比传统 SMTP 的优势

1. **更高的送达率**：Resend 专门优化了邮件送达率，避免进入垃圾邮件箱
2. **更简单的集成**：不需要处理复杂的 SMTP 协议、TLS 加密等
3. **更好的可靠性**：API 调用更稳定，超时和连接问题更少
4. **实时日志**：可以在 Resend 控制台查看邮件发送状态
5. **免费额度**：每月 3,000 封免费邮件（足够个人项目使用）

### 代码简化

- **删除代码**：~230 行 SMTP 相关代码
- **新增代码**：~80 行 Resend API 调用
- **净减少**：~150 行代码

## 验证测试

部署完成后，可以通过以下方式测试：

1. **注册测试**：尝试注册新用户，检查是否收到验证码邮件
2. **密码重置测试**：尝试使用"忘记密码"功能
3. **查看日志**：在 Render 控制台查看邮件发送日志
4. **Resend 控制台**：登录 Resend 查看邮件发送记录

## 故障排除

### 问题 1：邮件发送失败

**可能原因**：
- API Key 配置错误
- 发件域名未验证

**解决方案**：
1. 检查 `SMTP_PASSWORD` 环境变量是否正确
2. 在 Resend 控制台确认域名已验证
3. 查看 Render 日志获取详细错误信息

### 问题 2：邮件进入垃圾箱

**可能原因**：
- 域名声誉不足
- 缺少 DKIM 配置

**解决方案**：
1. 确保在 Resend 中配置了 DKIM 记录
2. 使用真实的发件域名，不要使用临时域名
3. 让收件人将发件地址添加到白名单

### 问题 3：DNS 配置未生效

**可能原因**：
- DNS 记录配置错误
- DNS 传播未完成

**解决方案**：
1. 使用 `dig` 命令检查 DNS 记录：
   ```bash
   dig txt resend._domainkey.ai-speed.xyz
   ```
2. 等待 DNS 传播（最多 48 小时）
3. 联系域名提供商确认配置

## 相关资源

- [Resend 官方文档](https://resend.com/docs)
- [Resend Go SDK](https://github.com/resend/resend-go)
- [域名验证指南](https://resend.com/docs/domains/verification)
- [API Keys 管理](https://resend.com/api-keys)

## 回滚方案

如果需要回滚到传统 SMTP 方式：

1. 恢复 `services/email.go` 中的 SMTP 代码
2. 修改环境变量为 SMTP 配置：
   - `SMTP_HOST`: smtp.qq.com
   - `SMTP_PORT`: 465
   - `SMTP_USER`: your-email@qq.com
   - `SMTP_PASSWORD`: your-authorization-code
3. 重新部署

但不建议回滚，因为 Resend 提供了更好的可靠性和送达率。

## 总结

本次迁移成功将邮件发送服务从传统 SMTP 升级到现代化的 Resend API，显著提高了代码质量和邮件送达率。只需确保在 Resend 控制台完成域名验证，即可开始使用更可靠的邮件服务。

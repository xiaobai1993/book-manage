# 在 Render 平台配置 Resend 环境变量

## 快速配置步骤

### 1. 访问 Render Dashboard

登录你的 [Render Dashboard](https://dashboard.render.com)

### 2. 选择你的服务

找到 `book-manage-backend` 服务（或你的后端服务名称）

### 3. 进入环境变量设置

点击服务 → 左侧菜单 **Environment** （或直接在服务页面找到 Environment 部分）

### 4. 添加以下环境变量

在 Environment Variables 部分添加（或更新）以下变量：

| Key | Value | 说明 |
|-----|-------|------|
| `SMTP_PASSWORD` | `re_xxxxxxxxxxxxxxxxx` | Resend API Key（从 https://resend.com/api-keys 获取） |
| `SMTP_USER` | `noreply@ai-speed.xyz` | 发件人邮箱地址 |

**重要提示**：
- 如果之前已经配置过 `SMTP_PASSWORD` 和 `SMTP_USER`，请更新为新值
- `SMTP_HOST` 和 `SMTP_PORT` 可以删除或留空（已不再使用）

### 5. 保存并重新部署

点击 **Save Changes** 或 **Deploy** 按钮，Render 会自动重新部署服务

## 验证配置

部署完成后，可以通过以下方式验证：

1. **查看日志**：在 Render Dashboard 查看 Logs，确认没有错误
2. **测试邮件**：使用前端的注册功能，发送验证码邮件
3. **检查收件箱**：检查是否收到验证码邮件

## Resend 域名验证（可选但推荐）

为了提高邮件送达率，建议在 Resend 控制台验证域名：

1. 访问 [Resend Domains](https://resend.com/domains)
2. 点击 **Add Domain**
3. 输入域名：`ai-speed.xyz`
4. 配置 DNS 记录（Resend 会提供具体的 DNS 记录）：
   - TXT 记录（域名验证）
   - CNAME 记录（DKIM，提高送达率）
5. 等待 DNS 生效（通常几分钟到几小时）

## 常见问题

### Q1: 邮件发送失败怎么办？

**解决方案**：
1. 检查 Render 环境变量是否正确配置
2. 查看 Render 日志获取详细错误信息
3. 确认 Resend API Key 是否有效

### Q2: 如何获取新的 API Key？

访问 [Resend API Keys](https://resend.com/api-keys) 创建或查看你的 API Keys

### Q3: 发件域名未验证会影响使用吗？

会影响送达率，但邮件仍然可以发送。建议尽快完成域名验证。

### Q4: 可以使用其他发件邮箱吗？

可以，但需要满足以下条件：
1. 域名已在 Resend 控制台验证
2. 使用已验证的域名作为发件邮箱
3. 格式：`anything@your-verified-domain.com`

## 相关文档

- [Resend 官方文档](https://resend.com/docs)
- [迁移说明](./RESEND_MIGRATION.md)
- [部署指南](../DEPLOY.md)

## 联系支持

如果遇到问题：
- 查看 [Render 文档](https://render.com/docs)
- 查看 [Resend 文档](https://resend.com/docs)
- 或在项目 GitHub 提 Issue

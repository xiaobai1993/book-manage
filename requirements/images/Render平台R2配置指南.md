# Render 平台 Cloudflare R2 配置指南

## 📋 概述

在 Render 平台部署时，Cloudflare R2 的配置需要通过**环境变量**来设置。与本地开发环境使用 YAML 配置文件不同，Render 使用环境变量来管理配置。

## 🔧 配置方式

### 方式一：通过 Render Dashboard 配置（推荐）

1. **登录 Render Dashboard**
   - 访问：https://dashboard.render.com/
   - 登录你的账户

2. **进入服务设置**
   - 找到你的服务（如 `book-manage-backend`）
   - 点击服务名称进入详情页

3. **打开环境变量设置**
   - 在左侧菜单找到 **Environment**（环境变量）
   - 或直接点击 **Environment Variables** 标签

4. **添加 R2 配置变量**
   点击 **Add Environment Variable**（添加环境变量），逐个添加以下变量：

   | Key（变量名） | Value（值） | 说明 |
   |--------------|-----------|------|
   | `R2_ACCOUNT_ID` | `your-r2-account-id` | Account ID（从S3端点URL提取） |
   | `R2_ACCESS_KEY_ID` | `your-r2-access-key-id` | S3 Access Key ID |
   | `R2_SECRET_ACCESS_KEY` | `your-r2-secret-access-key` | S3 Secret Access Key（敏感信息） |
   | `R2_BUCKET_NAME` | `your-bucket-name` | 存储桶名称 |
   | `R2_PUBLIC_URL` | `https://your-public-url.r2.dev` | 公开访问URL |
   | `R2_ENDPOINT` | `https://your-account-id.r2.cloudflarestorage.com` | S3端点URL |
   | `R2_REGION` | `auto` | 区域（默认auto） |

5. **保存配置**
   - 添加完所有变量后，点击 **Save Changes**（保存更改）
   - Render 会自动重新部署服务

### 方式二：通过 render.yaml 配置（已更新）

我已经更新了 `render.yaml` 文件，添加了 R2 配置。但需要注意：

**⚠️ 重要**：
- `R2_SECRET_ACCESS_KEY` 在 `render.yaml` 中设置为 `sync: false`
- 这意味着这个敏感信息**不会自动同步**到 Render
- 你需要在 Render Dashboard 中**手动设置**这个变量

**原因**：出于安全考虑，敏感信息（如密钥）不应该直接写在配置文件中。

### 方式三：使用 Render CLI（可选）

如果你使用 Render CLI，可以通过命令行设置：

```bash
# 安装 Render CLI
npm install -g render-cli

# 登录
render login

# 设置环境变量
render env:set R2_ACCOUNT_ID="your-r2-account-id" --service book-manage-backend
render env:set R2_ACCESS_KEY_ID="your-r2-access-key-id" --service book-manage-backend
render env:set R2_SECRET_ACCESS_KEY="your-r2-secret-access-key" --service book-manage-backend
render env:set R2_BUCKET_NAME="your-bucket-name" --service book-manage-backend
render env:set R2_PUBLIC_URL="https://your-public-url.r2.dev" --service book-manage-backend
render env:set R2_ENDPOINT="https://your-account-id.r2.cloudflarestorage.com" --service book-manage-backend
render env:set R2_REGION="auto" --service book-manage-backend
```

---

## 📝 详细配置步骤（Render Dashboard）

### 步骤 1：进入环境变量页面

1. 登录 Render Dashboard
2. 点击你的服务（如 `book-manage-backend`）
3. 在左侧菜单点击 **Environment**（环境变量）

### 步骤 2：添加环境变量

点击 **Add Environment Variable**（添加环境变量）按钮，逐个添加：

#### 变量 1：R2_ACCOUNT_ID
- **Key**: `R2_ACCOUNT_ID`
- **Value**: `your-r2-account-id`（请替换为实际值）
- 点击 **Save**

#### 变量 2：R2_ACCESS_KEY_ID
- **Key**: `R2_ACCESS_KEY_ID`
- **Value**: `your-r2-access-key-id`（请替换为实际值）
- 点击 **Save**

#### 变量 3：R2_SECRET_ACCESS_KEY（敏感信息）
- **Key**: `R2_SECRET_ACCESS_KEY`
- **Value**: `your-r2-secret-access-key`（请替换为实际值）
- **⚠️ 注意**：这是敏感信息，确保不要泄露
- 点击 **Save**

#### 变量 4：R2_BUCKET_NAME
- **Key**: `R2_BUCKET_NAME`
- **Value**: `your-bucket-name`（请替换为实际值）
- 点击 **Save**

#### 变量 5：R2_PUBLIC_URL
- **Key**: `R2_PUBLIC_URL`
- **Value**: `https://your-public-url.r2.dev`（请替换为实际值）
- 点击 **Save**

#### 变量 6：R2_ENDPOINT
- **Key**: `R2_ENDPOINT`
- **Value**: `https://your-account-id.r2.cloudflarestorage.com`（请替换为实际值）
- 点击 **Save**

#### 变量 7：R2_REGION
- **Key**: `R2_REGION`
- **Value**: `auto`
- 点击 **Save**

### 步骤 3：验证配置

添加完所有变量后，你应该能看到 7 个环境变量：

```
R2_ACCOUNT_ID = your-r2-account-id
R2_ACCESS_KEY_ID = your-r2-access-key-id
R2_SECRET_ACCESS_KEY = your-r2-secret-access-key
R2_BUCKET_NAME = your-bucket-name
R2_PUBLIC_URL = https://your-public-url.r2.dev
R2_ENDPOINT = https://your-account-id.r2.cloudflarestorage.com
R2_REGION = auto
```

### 步骤 4：重新部署

配置完成后，Render 会自动触发重新部署。如果没有自动部署：

1. 点击 **Manual Deploy**（手动部署）
2. 选择 **Deploy latest commit**（部署最新提交）
3. 等待部署完成

---

## 🔍 配置验证

### 方法一：查看部署日志

1. 在 Render Dashboard 中，点击 **Logs**（日志）标签
2. 查看启动日志，应该能看到：
   - 如果配置正确：R2 服务初始化成功（或没有错误）
   - 如果配置错误：会显示 R2 服务初始化失败的错误信息

### 方法二：测试图片上传

1. 部署完成后，测试图片上传功能
2. 如果上传成功，说明配置正确
3. 如果失败，查看错误信息并检查配置

---

## ⚠️ 常见问题

### Q1: 环境变量设置后没有生效？

**A**: 
- 确保点击了 **Save Changes**（保存更改）
- 检查服务是否已重新部署
- 查看部署日志确认环境变量已加载

### Q2: 如何知道环境变量是否正确设置？

**A**: 
- 在 Render Dashboard 的 Environment 页面可以看到所有环境变量
- 注意：`R2_SECRET_ACCESS_KEY` 的值会显示为 `••••••`（隐藏），这是正常的

### Q3: 可以批量导入环境变量吗？

**A**: 
- Render Dashboard 支持批量导入
- 点击 **Import from file**（从文件导入）
- 格式：每行一个变量，格式为 `KEY=VALUE`

### Q4: 环境变量区分大小写吗？

**A**: 
- 是的，环境变量名区分大小写
- 确保使用正确的大小写：`R2_ACCOUNT_ID`（不是 `r2_account_id`）

### Q5: 如果忘记设置某个环境变量会怎样？

**A**: 
- 如果缺少必需的环境变量，R2 服务会初始化失败
- 图片上传功能将不可用
- 但不会影响其他功能（如图书管理、借阅等）
- 查看日志可以看到具体的错误信息

---

## 🔐 安全建议

1. **保护 Secret Access Key**
   - 不要在公开场合分享 `R2_SECRET_ACCESS_KEY`
   - 不要提交到代码仓库
   - 在 Render Dashboard 中，这个值会显示为 `••••••`（隐藏）

2. **定期轮换密钥**
   - 建议每 6-12 个月更换一次 R2 API Token
   - 更换后，更新 Render 中的环境变量

3. **使用环境变量而不是配置文件**
   - 生产环境使用环境变量（更安全）
   - 开发环境可以使用配置文件（方便）

---

## 📊 配置对比

### 本地开发环境（config/env.yaml）
```yaml
cloudflare_r2:
  account_id: "your-r2-account-id"
  access_key_id: "your-r2-access-key-id"
  secret_access_key: "your-r2-secret-access-key"
  bucket_name: "your-bucket-name"
  public_url: "https://your-public-url.r2.dev"
  endpoint: "https://your-account-id.r2.cloudflarestorage.com"
  region: "auto"
```

### Render 生产环境（环境变量）
```bash
R2_ACCOUNT_ID=your-r2-account-id
R2_ACCESS_KEY_ID=your-r2-access-key-id
R2_SECRET_ACCESS_KEY=your-r2-secret-access-key
R2_BUCKET_NAME=your-bucket-name
R2_PUBLIC_URL=https://your-public-url.r2.dev
R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
R2_REGION=auto
```

**代码会自动识别**：无论是从配置文件还是环境变量读取，代码都能正确处理。

---

## ✅ 配置检查清单

完成配置后，确认：

- [ ] 所有 7 个环境变量都已添加
- [ ] `R2_SECRET_ACCESS_KEY` 已正确设置（值会隐藏显示）
- [ ] 服务已重新部署
- [ ] 查看日志确认没有 R2 初始化错误
- [ ] 测试图片上传功能（可选）

---

## 🎉 完成

配置完成后，你的图书管理系统就可以使用 Cloudflare R2 存储图片了！

如果遇到问题，可以：
1. 查看 Render 部署日志
2. 检查环境变量是否正确设置
3. 测试图片上传功能


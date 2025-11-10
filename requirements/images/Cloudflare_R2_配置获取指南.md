# Cloudflare R2 配置信息获取指南

## 📋 需要配置的信息清单

你需要准备以下 6 个配置项：

| 配置项 | 说明 | 是否必需 | 示例值 |
|--------|------|---------|--------|
| `account_id` | Cloudflare 账户 ID | ✅ 必需 | `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6` |
| `access_key_id` | R2 API Token 的访问密钥 ID | ✅ 必需 | `your-access-key-id` |
| `secret_access_key` | R2 API Token 的密钥 | ✅ 必需 | `your-secret-access-key` |
| `bucket_name` | R2 存储桶名称 | ✅ 必需 | `book-covers-prod` |
| `public_url` | 图片公开访问 URL | ✅ 必需 | `https://pub-xxxxx.r2.dev` |
| `region` | 存储区域 | ⚠️ 可选 | `auto`（默认值） |

---

## 🔍 详细获取步骤

### 步骤 1：获取 Account ID（账户 ID）

#### 方法一：从 Dashboard 首页获取
1. 登录 Cloudflare Dashboard：https://dash.cloudflare.com/
2. 在右侧边栏找到 **Account ID**
3. 点击复制按钮复制 Account ID
4. 格式示例：`a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`（32 位字符）

#### 方法二：从任意域名页面获取
1. 选择任意一个域名
2. 在右侧边栏可以看到 **Account ID**
3. 复制该 ID

**⚠️ 注意**：Account ID 是全局的，一个账户只有一个 Account ID。

---

### 步骤 2：创建 R2 存储桶并获取 Bucket Name

1. 登录 Cloudflare Dashboard：https://dash.cloudflare.com/
2. 在左侧菜单栏找到 **R2**，点击进入
3. 如果是第一次使用，可能需要：
   - 点击 **Get started** 或 **Create bucket**
   - 阅读并同意服务条款
4. 点击 **Create bucket**（创建存储桶）
5. 填写存储桶信息：
   - **Bucket name**（存储桶名称）：
     - 开发环境：`book-covers-dev`
     - 生产环境：`book-covers-prod`
   - **Location**（位置）：选择离你最近的区域
     - `apac` - 亚太地区（推荐中国用户）
     - `wnam` - 美国西部
     - `enam` - 美国东部
     - `eeur` - 欧洲东部
     - `weur` - 欧洲西部
6. 点击 **Create bucket** 完成创建
7. **记录存储桶名称**：这就是你需要的 `bucket_name`

**📝 提示**：
- 存储桶名称必须全局唯一
- 建议使用环境后缀区分：`book-covers-dev`、`book-covers-prod`
- 存储桶创建后可以随时重命名（在设置中）

---

### 步骤 3：创建 R2 API Token 并获取 Access Key 和 Secret

#### 3.1 创建 API Token

1. 在 R2 页面，点击右上角的 **Manage R2 API Tokens**（管理 R2 API 令牌）
   - 或者直接访问：https://dash.cloudflare.com/?to=/:account/r2/api-tokens

2. 点击 **Create API token**（创建 API 令牌）

3. 填写 Token 配置：
   - **Token name**（令牌名称）：`book-manage-r2-token`（可自定义）
   - **Permissions**（权限）：选择 **Object Read & Write**（对象读写）
   - **TTL**（有效期）：
     - 开发环境：选择 **Never expire**（永不过期）
     - 生产环境：建议设置过期时间（如 1 年），到期后重新创建
   - **R2 Token Scopes**（R2 令牌作用域）：
     - 选择你刚创建的存储桶（如 `book-covers-dev`）
     - 或者选择 **All buckets**（所有存储桶）- 不推荐，权限过大

4. 点击 **Create API Token**（创建 API 令牌）

#### 3.2 获取 Access Key ID 和 Secret Access Key

**⚠️ 重要**：创建 Token 后，系统会显示一次性的密钥信息，请立即保存！

1. 创建成功后，会显示一个弹窗，包含：
   - **Access Key ID**：类似 `your-access-key-id-here`
   - **Secret Access Key**：类似 `your-secret-access-key-here`（较长）

2. **立即复制并保存**：
   - 点击 **Copy** 按钮复制每个密钥
   - 或者手动复制
   - **Secret Access Key 只显示一次**，关闭后无法再次查看！

3. **安全保存**：
   - 建议保存到密码管理器
   - 或保存到本地加密文件
   - **不要提交到代码仓库**

4. 点击 **I've copied the secret**（我已复制密钥）关闭弹窗

**📝 提示**：
- 如果忘记 Secret Access Key，需要删除旧 Token 并重新创建
- 可以创建多个 Token 用于不同环境（开发、生产）

---

### 步骤 4：配置公开访问并获取 Public URL

有两种方式配置公开访问：

#### 方式一：使用 R2.dev Subdomain（推荐，最简单）

1. 在 R2 页面，点击你创建的存储桶名称进入详情页

2. 点击 **Settings**（设置）标签

3. 找到 **Public Access**（公开访问）部分

4. 点击 **Connect Domain**（连接域名）或 **Enable R2.dev Subdomain**（启用 R2.dev 子域名）

5. 选择 **R2.dev Subdomain** 选项

6. 系统会自动分配一个子域名，格式如：
   - `https://pub-xxxxx.r2.dev`
   - 或 `https://pub-xxxxx.r2.dev`

7. **复制这个 URL**：这就是你需要的 `public_url`

8. 点击 **Save**（保存）

**📝 提示**：
- R2.dev 子域名是免费的
- 子域名是自动分配的，无法自定义
- 启用后，存储桶中的所有文件都可以通过这个 URL 公开访问

#### 方式二：使用自定义域名（可选，更专业）

1. 在存储桶设置中，选择 **Custom Domain**（自定义域名）

2. 输入你的域名（如 `images.yourdomain.com`）

3. 按照提示配置 DNS 记录（CNAME）

4. 等待 DNS 生效后，使用自定义域名作为 `public_url`

**⚠️ 注意**：自定义域名需要你有自己的域名并配置 DNS，对于免费方案，推荐使用 R2.dev 子域名。

---

### 步骤 5：Region（区域）- 可选配置

- **默认值**：`auto`（自动）
- **说明**：Cloudflare 会自动选择最佳区域
- **可选值**：
  - `apac` - 亚太地区
  - `wnam` - 美国西部
  - `enam` - 美国东部
  - `eeur` - 欧洲东部
  - `weur` - 欧洲西部

**📝 建议**：使用默认值 `auto` 即可，除非有特殊需求。

---

## 📝 配置信息汇总表

获取完所有信息后，填写以下表格：

```
Account ID:        _____________________________
Access Key ID:     _____________________________
Secret Access Key: _____________________________
Bucket Name:       _____________________________
Public URL:        _____________________________
Region:            auto（默认值）
```

---

## 🔐 配置到项目中的方式

### 方式一：开发环境 - 配置文件（config/env.yaml）

在 `config/env.yaml` 中添加：

```yaml
cloudflare_r2:
  account_id: "你的Account ID"
  access_key_id: "你的Access Key ID"
  secret_access_key: "你的Secret Access Key"
  bucket_name: "book-covers-dev"
  public_url: "https://pub-xxxxx.r2.dev"
  region: "auto"
```

### 方式二：生产环境 - 环境变量（推荐，更安全）

在生产环境（如 Render、Vercel 等）中设置环境变量：

```bash
R2_ACCOUNT_ID=你的Account ID
R2_ACCESS_KEY_ID=你的Access Key ID
R2_SECRET_ACCESS_KEY=你的Secret Access Key
R2_BUCKET_NAME=book-covers-prod
R2_PUBLIC_URL=https://pub-xxxxx.r2.dev
R2_REGION=auto
```

---

## ✅ 验证配置是否正确

配置完成后，可以通过以下方式验证：

1. **测试上传**：尝试上传一张测试图片
2. **检查 URL**：确认返回的图片 URL 可以正常访问
3. **查看存储桶**：在 Cloudflare Dashboard 的 R2 页面查看文件是否已上传

---

## 🆘 常见问题

### Q1: 找不到 Account ID？
**A**: Account ID 在 Dashboard 右侧边栏，如果看不到，尝试：
- 刷新页面
- 选择任意一个域名，右侧会显示 Account ID

### Q2: Secret Access Key 忘记保存了？
**A**: 需要删除旧 Token 并重新创建。步骤：
1. 进入 R2 API Tokens 页面
2. 找到对应的 Token
3. 点击删除
4. 重新创建并保存密钥

### Q3: R2.dev 子域名无法访问？
**A**: 检查：
1. 是否已启用 R2.dev Subdomain
2. 存储桶的 Public Access 设置是否正确
3. 文件路径是否正确（格式：`{public_url}/book-covers/{filename}`）

### Q4: 如何查看已创建的 Token？
**A**: 
1. 进入 R2 API Tokens 页面
2. 可以看到所有已创建的 Token 列表
3. 但 Secret Access Key 不会再次显示（安全考虑）

### Q5: 存储桶名称可以修改吗？
**A**: 可以，在存储桶的 Settings 中可以重命名，但建议创建时就使用正确的名称。

---

## 📚 参考链接

- Cloudflare Dashboard：https://dash.cloudflare.com/
- R2 文档：https://developers.cloudflare.com/r2/
- R2 API Tokens 管理：https://dash.cloudflare.com/?to=/:account/r2/api-tokens

---

## ⚠️ 安全提醒

1. **不要将 Secret Access Key 提交到代码仓库**
2. **使用环境变量存储敏感信息**
3. **定期轮换 API Token**（建议每 6-12 个月）
4. **限制 Token 权限范围**（只给必要的存储桶权限）
5. **生产环境和开发环境使用不同的 Token**


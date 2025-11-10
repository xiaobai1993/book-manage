# Cloudflare R2 新手详细操作指南

## 📖 目录
1. [注册 Cloudflare 账户](#1-注册-cloudflare-账户)
2. [创建 R2 存储桶](#2-创建-r2-存储桶)
3. [创建 API Token](#3-创建-api-token)
4. [配置公开访问](#4-配置公开访问)
5. [获取所有配置信息](#5-获取所有配置信息)
6. [验证配置](#6-验证配置)

---

## 1. 注册 Cloudflare 账户

### 步骤 1.1：访问注册页面
1. 打开浏览器，访问：https://dash.cloudflare.com/sign-up
2. 如果已有账户，点击右上角 **Log in**（登录）

### 步骤 1.2：填写注册信息
1. **Email**（邮箱）：输入你的邮箱地址
2. **Password**（密码）：设置密码（至少 8 位）
3. 勾选同意服务条款
4. 点击 **Sign up**（注册）

### 步骤 1.3：验证邮箱
1. 检查邮箱收件箱
2. 找到 Cloudflare 发送的验证邮件
3. 点击邮件中的验证链接
4. 完成邮箱验证

### 步骤 1.4：登录 Dashboard
1. 验证完成后，会自动跳转到 Dashboard
2. 或者访问：https://dash.cloudflare.com/
3. 使用邮箱和密码登录

**✅ 完成标志**：能看到 Cloudflare Dashboard 主页面

---

## 2. 创建 R2 存储桶

### ⚠️ 重要：关于支付方式要求

**Cloudflare R2 需要绑定支付方式才能使用**，但这是为了：
- 验证账户身份，防止滥用
- 超出免费额度时的计费（但免费额度很大，一般不会超）

**好消息**：
- R2 有大量免费额度（每月 10GB 存储，1000 万次读取）
- 对于图书管理系统，基本不会超出免费额度
- 绑定支付方式后，只要不超出免费额度，**不会产生任何费用**

### 💳 中国用户可用的支付方式

#### 方式一：国际信用卡（推荐）
支持以下类型的信用卡：
- **Visa**（维萨卡）- 最常用
- **Mastercard**（万事达卡）- 最常用
- **American Express**（美国运通）- 部分支持

**中国可用的银行卡**：
1. **中国银行**：Visa、Mastercard 双币信用卡
2. **工商银行**：Visa、Mastercard 双币信用卡
3. **建设银行**：Visa、Mastercard 双币信用卡
4. **招商银行**：Visa、Mastercard 双币信用卡
5. **中信银行**：Visa、Mastercard 双币信用卡
6. **浦发银行**：Visa、Mastercard 双币信用卡
7. **其他银行**：大部分银行都提供支持国际支付的信用卡

**如何申请**：
- 联系你的银行客服咨询"双币信用卡"或"国际信用卡"
- 通常需要提供收入证明等材料
- 申请周期：1-2 周

#### 方式二：虚拟信用卡（备选方案）
如果暂时没有国际信用卡，可以考虑：
- **NobePay**：提供虚拟 Visa 卡（需要充值）
- **Depay**：提供虚拟信用卡服务
- **其他虚拟卡服务**：搜索"虚拟信用卡"了解

**⚠️ 注意**：虚拟信用卡服务需要自行了解其可靠性和费用

#### 方式三：PayPal（如果支持）
- 部分情况下 Cloudflare 可能支持 PayPal
- 需要绑定 PayPal 账户
- PayPal 可以绑定中国银行卡（需要验证）

### 📝 绑定支付方式步骤

1. 当提示需要绑定支付方式时：
   - 点击 **Add Payment Method**（添加支付方式）
   - 选择 **Credit Card**（信用卡）

2. 填写信用卡信息：
   - **Card Number**（卡号）：输入信用卡号
   - **Expiry Date**（有效期）：MM/YY 格式
   - **CVV**（安全码）：卡片背面的 3 位数字
   - **Cardholder Name**（持卡人姓名）：与卡片一致
   - **Billing Address**（账单地址）：填写真实地址

3. 点击 **Submit**（提交）

4. Cloudflare 可能会：
   - 验证卡片（可能扣除 1 美元验证，稍后退还）
   - 或直接验证通过

### ⚠️ 安全提醒

1. **Cloudflare 是可信的公司**，可以安全绑定
2. **设置使用限额**（如果担心）：
   - 在 Cloudflare Dashboard 可以设置使用限额
   - 超出限额会自动停止服务，不会产生意外费用

3. **监控使用量**：
   - 定期检查 R2 使用量
   - 免费额度很大，一般不会超出

### 💡 免费额度说明

Cloudflare R2 免费额度：
- **存储空间**：每月 10GB（免费）
- **读取操作**：每月 1000 万次（免费）
- **写入操作**：每月 100 万次（免费）

**对于图书管理系统**：
- 假设 1000 本图书，每本封面 500KB
- 总存储：约 500MB（远低于 10GB）
- 每月访问：10 万次（远低于 1000 万次）

**结论**：完全在免费额度内，不会产生费用 ✅

### 步骤 2.1：找到 R2 服务入口
1. 登录后，在左侧菜单栏找到 **R2** 图标
   - 图标是一个类似云存储的图标
   - 如果没看到，可能需要向下滚动菜单
2. 点击 **R2** 进入 R2 服务页面

### 步骤 2.2：首次使用 R2（如果是第一次）
1. 如果看到 **Get started** 或 **Create bucket** 按钮，直接点击
2. 如果看到服务条款或介绍页面：
   - 阅读并了解 R2 服务
   - 点击 **I understand**（我理解）或 **Continue**（继续）
   - 可能需要同意服务条款
3. **如果提示需要绑定支付方式**：
   - 按照上面的"绑定支付方式步骤"操作
   - 使用支持国际支付的信用卡
   - 绑定完成后继续创建存储桶

### 步骤 2.3：创建存储桶
1. 在 R2 页面，点击右上角的 **Create bucket**（创建存储桶）按钮
   - 按钮通常是蓝色或绿色的
   - 位置在页面右上角

2. 填写存储桶信息：
   - **Bucket name**（存储桶名称）：
     - 输入：`book-covers-dev`（开发环境）
     - 或：`book-covers-prod`（生产环境）
     - **注意**：名称必须全局唯一，如果提示已存在，可以加后缀，如 `book-covers-dev-2025`
   
   - **Location**（位置/区域）：
     - 点击下拉菜单
     - 选择 **apac**（亚太地区）- 推荐中国用户选择
     - 其他选项：
       - `wnam` - 美国西部
       - `enam` - 美国东部
       - `eeur` - 欧洲东部
       - `weur` - 欧洲西部

3. 点击 **Create bucket**（创建存储桶）按钮

### 步骤 2.4：确认创建成功
1. 创建成功后，会看到存储桶列表
2. 你刚创建的存储桶会显示在列表中
3. **记录存储桶名称**：这就是你需要的 `bucket_name`（如：`book-covers-dev`）

**✅ 完成标志**：能在 R2 页面看到你创建的存储桶

**📝 提示**：
- 存储桶名称可以随时修改（在设置中）
- 可以创建多个存储桶用于不同环境
- 存储桶名称在 Cloudflare 全局必须唯一

---

## 3. 创建 API Token

### 步骤 3.1：进入 API Token 管理页面
有两种方式：

**方式一（推荐）**：
1. 在 R2 页面，点击右上角的 **Manage R2 API Tokens**（管理 R2 API 令牌）
   - 按钮通常在页面右上角，可能在 **Create bucket** 按钮旁边

**方式二（直接访问）**：
1. 直接在浏览器地址栏输入：
   ```
   https://dash.cloudflare.com/?to=/:account/r2/api-tokens
   ```
2. 按回车访问

### 步骤 3.2：创建新的 API Token
1. 在 API Tokens 页面，点击 **Create API token**（创建 API 令牌）按钮
   - 按钮通常是蓝色或绿色的，在页面顶部或右上角

2. 填写 Token 配置信息：

   **a) Token name（令牌名称）**：
   - 输入：`book-manage-r2-token`
   - 或自定义名称，用于识别这个 Token 的用途
   - 建议包含环境信息，如：`book-manage-dev`、`book-manage-prod`

   **b) Permissions（权限）**：
   - 在下拉菜单中选择：**Object Read & Write**（对象读写）
   - 这个权限允许上传、下载、删除文件
   - **不要选择** "Admin" 权限（权限过大）

   **c) TTL（有效期）**：
   - 开发环境：选择 **Never expire**（永不过期）
   - 生产环境：建议选择 **Custom**（自定义），设置 1 年或更长时间
   - 点击日期选择器设置过期时间

   **d) R2 Token Scopes（R2 令牌作用域）**：
   - 选择 **Specific bucket**（特定存储桶）
   - 在下拉菜单中选择你刚创建的存储桶（如 `book-covers-dev`）
   - **不要选择** "All buckets"（所有存储桶）- 权限过大，不安全

3. 点击页面底部的 **Create API Token**（创建 API 令牌）按钮

### 步骤 3.3：保存 Access Key 和 Secret Key（⚠️ 非常重要）
1. 创建成功后，会弹出一个对话框或新页面
2. 显示两个重要信息：

   **Access Key ID**：
   - 类似：`your-access-key-id-here`
   - 点击旁边的 **Copy**（复制）按钮复制
   - 或手动选中文本复制（Ctrl+C / Cmd+C）

   **Secret Access Key**：
   - 类似：`your-secret-access-key-here`（较长的一串字符）
   - 点击旁边的 **Copy**（复制）按钮复制
   - **⚠️ 重要**：这个密钥只显示一次，关闭后无法再次查看！

3. **立即保存这两个密钥**：
   - 建议保存到文本文件（加密保存）
   - 或保存到密码管理器
   - 或暂时保存到本地文档
   - **绝对不要**提交到代码仓库

4. 确认已复制后，点击 **I've copied the secret**（我已复制密钥）按钮
   - 或点击 **Close**（关闭）按钮

**✅ 完成标志**：已保存 Access Key ID 和 Secret Access Key

**📝 提示**：
- 如果忘记保存 Secret Access Key，需要删除这个 Token 并重新创建
- 可以创建多个 Token 用于不同环境
- Token 创建后可以在列表中看到，但 Secret 不会再次显示

---

## 4. 配置公开访问

### 步骤 4.1：进入存储桶设置
1. 回到 R2 主页面（点击左侧菜单的 **R2**）
2. 在存储桶列表中，点击你创建的存储桶名称（如 `book-covers-dev`）
3. 进入存储桶详情页面

### 步骤 4.2：打开设置页面
1. 在存储桶详情页面，点击顶部的 **Settings**（设置）标签
   - 标签通常在页面顶部，可能在 **Overview**（概览）旁边

### 步骤 4.3：启用 R2.dev Subdomain
1. 在设置页面，向下滚动找到 **Public Access**（公开访问）部分
2. 点击 **Connect Domain**（连接域名）按钮
   - 或点击 **Enable R2.dev Subdomain**（启用 R2.dev 子域名）

3. 选择 **R2.dev Subdomain** 选项
   - 这是最简单的方式，免费且无需配置 DNS

4. 系统会自动分配一个子域名，格式类似：
   ```
   https://pub-xxxxx.r2.dev
   ```
   - `xxxxx` 是系统自动生成的随机字符

5. **复制这个 URL**：
   - 选中整个 URL（包括 `https://`）
   - 复制（Ctrl+C / Cmd+C）
   - 这就是你需要的 `public_url`

6. 点击 **Save**（保存）或 **Enable**（启用）按钮

### 步骤 4.4：确认配置成功
1. 保存后，**Public Access** 部分会显示已启用的子域名
2. 状态应该显示为 **Enabled**（已启用）
3. **再次确认并复制 Public URL**，确保正确

**✅ 完成标志**：能看到已启用的 R2.dev 子域名 URL

**📝 提示**：
- R2.dev 子域名是免费的
- 子域名是自动分配的，无法自定义
- 启用后，存储桶中的所有文件都可以通过这个 URL 公开访问
- 如果不想使用公开访问，可以稍后禁用

---

## 5. 获取所有配置信息

### 步骤 5.1：获取 Account ID
1. 在 Cloudflare Dashboard 任意页面
2. 查看右侧边栏（如果没有显示，可能需要点击某个域名）
3. 找到 **Account ID** 字段
4. 点击旁边的复制图标（📋）或直接复制文本
5. 格式：32 位字符，如 `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`

**如果找不到 Account ID**：
- 方法 1：点击左侧菜单的任意一个域名，右侧会显示 Account ID
- 方法 2：在 Dashboard 首页，右侧边栏通常有 Account ID
- 方法 3：在 R2 页面的 URL 中可以看到，格式：`/account/{account-id}/r2/`

### 步骤 5.2：汇总所有配置信息
现在你应该已经获取了所有配置信息，填写下面的表格：

```
✅ Account ID:        _____________________________
   （从 Dashboard 右侧边栏获取）

✅ Access Key ID:     _____________________________
   （从步骤 3.3 中保存的）

✅ Secret Access Key: _____________________________
   （从步骤 3.3 中保存的，只显示一次）

✅ Bucket Name:      _____________________________
   （从步骤 2.3 中创建的，如：book-covers-dev）

✅ Public URL:        _____________________________
   （从步骤 4.3 中获取的 R2.dev 子域名）

✅ Region:           auto
   （使用默认值 auto 即可）
```

---

## 6. 验证配置

### 步骤 6.1：检查配置完整性
确认你已经获取了所有 6 个配置项：
- [ ] Account ID
- [ ] Access Key ID
- [ ] Secret Access Key
- [ ] Bucket Name
- [ ] Public URL
- [ ] Region（使用 `auto`）

### 步骤 6.2：测试存储桶访问
1. 在 R2 页面，点击你的存储桶
2. 确认能看到存储桶详情页面
3. 存储桶应该是空的（还没有上传文件）

### 步骤 6.3：测试 Public URL（可选）
1. 在存储桶中，尝试上传一个测试文件：
   - 点击 **Upload**（上传）按钮
   - 选择一个测试图片
   - 上传成功后，点击文件名
   - 查看是否能通过 Public URL 访问
   - URL 格式：`{public_url}/文件名`

2. 如果能看到图片，说明配置正确 ✅

---

## 📝 配置信息使用方式

### 开发环境配置（config/env.yaml）
将配置信息添加到 `config/env.yaml`：

```yaml
cloudflare_r2:
  account_id: "粘贴你的 Account ID"
  access_key_id: "粘贴你的 Access Key ID"
  secret_access_key: "粘贴你的 Secret Access Key"
  bucket_name: "book-covers-dev"
  public_url: "粘贴你的 Public URL（如：https://pub-xxxxx.r2.dev）"
  region: "auto"
```

### 生产环境配置（环境变量）
在生产环境平台（如 Render、Vercel）设置环境变量：

```bash
R2_ACCOUNT_ID=你的Account ID
R2_ACCESS_KEY_ID=你的Access Key ID
R2_SECRET_ACCESS_KEY=你的Secret Access Key
R2_BUCKET_NAME=book-covers-prod
R2_PUBLIC_URL=你的Public URL
R2_REGION=auto
```

---

## 🆘 常见问题解决

### Q1: 找不到 R2 菜单？
**A**: 
- 确保已登录 Cloudflare Dashboard
- 检查左侧菜单是否折叠，尝试展开
- 如果还是没有，可能需要先添加一个域名（免费账户也可以添加域名）

### Q2: 创建存储桶时提示名称已存在？
**A**: 
- 存储桶名称必须全局唯一
- 尝试添加后缀，如：`book-covers-dev-2025`、`book-covers-dev-yourname`

### Q3: 找不到 "Manage R2 API Tokens" 按钮？
**A**: 
- 尝试直接访问：https://dash.cloudflare.com/?to=/:account/r2/api-tokens
- 或者在 R2 页面右上角查找类似 "API" 或 "Tokens" 的链接

### Q4: Secret Access Key 忘记保存了？
**A**: 
- 需要删除旧 Token 并重新创建
- 在 API Tokens 页面，找到对应的 Token，点击删除
- 然后重新创建并立即保存密钥

### Q5: 找不到 Public Access 设置？
**A**: 
- 确保已进入存储桶详情页面（点击存储桶名称）
- 点击顶部的 **Settings**（设置）标签
- 向下滚动找到 **Public Access** 部分

### Q6: R2.dev 子域名无法访问图片？
**A**: 
- 确认已启用 R2.dev Subdomain
- 检查文件路径是否正确
- 尝试在浏览器直接访问：`{public_url}/文件名`

### Q7: 找不到 Account ID？
**A**: 
- 点击左侧菜单的任意域名，右侧会显示 Account ID
- 或者在 Dashboard 首页右侧边栏查找
- 也可以在浏览器地址栏查看 URL，格式：`/account/{account-id}/`

---

## ✅ 完成检查清单

完成所有步骤后，确认：

- [ ] 已注册 Cloudflare 账户并登录
- [ ] 已创建 R2 存储桶
- [ ] 已创建 API Token 并保存了 Access Key 和 Secret Key
- [ ] 已启用 R2.dev Subdomain 并获取了 Public URL
- [ ] 已获取 Account ID
- [ ] 已将所有配置信息保存到安全的地方
- [ ] 已测试存储桶可以正常访问

---

## 🎉 下一步

配置信息获取完成后，你可以：
1. 将这些配置信息添加到项目配置文件中
2. 按照《Cloudflare_R2_实施方案.md》文档实现代码
3. 测试图片上传功能

如果遇到任何问题，随时告诉我！我会帮你解决。


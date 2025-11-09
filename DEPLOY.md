# 部署指南

本指南将帮助你将图书管理系统部署到线上环境，使用以下技术栈：
- **前端**: Vercel
- **后端**: Render
- **数据库**: Supabase (PostgreSQL)

## 📋 前置准备

1. **GitHub 账号**：确保代码已推送到 GitHub
2. **Supabase 账号**：注册 [Supabase](https://supabase.com)
3. **Vercel 账号**：注册 [Vercel](https://vercel.com)
4. **Render 账号**：注册 [Render](https://render.com)

---

## 🗄️ 第一步：设置 Supabase 数据库

### 1.1 创建 Supabase 项目

1. 登录 [Supabase Dashboard](https://app.supabase.com)
2. 点击 "New Project"
3. 填写项目信息：
   - **Name**: book-manage（或你喜欢的名字）
   - **Database Password**: 设置一个强密码（**请务必保存！**）
   - **Region**: 选择离你最近的区域（如 `Southeast Asia (Singapore)`）
4. 点击 "Create new project"，等待创建完成（约 2 分钟）

### 1.2 导入数据库结构

1. 在 Supabase Dashboard 中，点击左侧菜单的 **SQL Editor**
2. 点击 "New query"
3. 打开项目中的 `data_postgresql.sql` 文件，复制全部内容
4. 粘贴到 SQL Editor 中
5. 点击 "Run" 执行 SQL 脚本
6. 确认表创建成功（应该看到 `user`、`book`、`borrow_record` 三个表）

### 1.3 获取数据库连接信息

1. 在 Supabase Dashboard 中，点击左侧菜单的 **Settings** → **Database**
2. 找到 **Connection string** 部分
3. 选择 **URI** 格式，复制连接字符串（格式类似：`postgresql://postgres:[YOUR-PASSWORD]@db.xxx.supabase.co:5432/postgres`）
4. **保存这个连接字符串**，后续在 Render 中会用到

---

## 🚀 第二步：部署后端到 Render

### 2.1 创建 Web Service

1. 登录 [Render Dashboard](https://dashboard.render.com)
2. 点击 "New +" → "Web Service"
3. 连接你的 GitHub 仓库
4. 选择 `book-manage` 仓库

### 2.2 配置服务

填写以下信息：

- **Name**: `book-manage-backend`（或你喜欢的名字）
- **Environment**: `Go`
- **Region**: 选择离你最近的区域
- **Branch**: `master`（或你的主分支）
- **Root Directory**: 留空（根目录）
- **Build Command**: 
  ```bash
  go mod download && go build -o book-manage
  ```
- **Start Command**: 
  ```bash
  ./book-manage
  ```

### 2.3 配置环境变量

在 **Environment Variables** 部分，添加以下变量：

| 变量名 | 值 | 说明 |
|--------|-----|------|
| `DB_TYPE` | `postgres` | 数据库类型 |
| `DATABASE_URL` | `你的 Supabase 连接字符串` | 从 Supabase 复制的完整连接字符串 |
| `JWT_SECRET` | `生成一个随机字符串` | 用于 JWT 加密（可以使用在线工具生成） |
| `PORT` | `8080` | 服务器端口（Render 会自动设置，但可以显式指定） |
| `ADMIN_EMAILS` | `admin@lib.com` | 管理员邮箱（多个用逗号分隔） |
| `SMTP_HOST` | `smtp.qq.com` | 邮箱 SMTP 主机（可选） |
| `SMTP_PORT` | `465` | 邮箱 SMTP 端口（可选） |
| `SMTP_USER` | `你的邮箱` | 邮箱用户名（可选） |
| `SMTP_PASSWORD` | `你的邮箱密码` | 邮箱密码（可选） |

**重要提示**：
- `DATABASE_URL` 应该包含密码，格式：`postgresql://postgres:YOUR_PASSWORD@db.xxx.supabase.co:5432/postgres`
- `JWT_SECRET` 建议使用强随机字符串，可以使用：`openssl rand -base64 32` 生成

### 2.4 部署

1. 点击 "Create Web Service"
2. Render 会自动开始构建和部署
3. 等待部署完成（约 5-10 分钟）
4. 部署成功后，你会得到一个 URL，例如：`https://book-manage-backend.onrender.com`
5. **保存这个 URL**，后续配置前端时会用到

### 2.5 测试后端

在浏览器中访问：`https://your-backend-url.onrender.com/api/user/login`（应该会返回错误，但说明服务已启动）

---

## 🎨 第三步：部署前端到 Vercel

### 3.1 创建 Vercel 项目

1. 登录 [Vercel Dashboard](https://vercel.com/dashboard)
2. 点击 "Add New..." → "Project"
3. 导入你的 GitHub 仓库 `book-manage`
4. 点击 "Import"

### 3.2 配置项目

在项目配置页面：

- **Framework Preset**: `Vite`
- **Root Directory**: `frontend`
- **Build Command**: `npm run build`
- **Output Directory**: `dist`
- **Install Command**: `npm install`

### 3.3 配置环境变量

在 **Environment Variables** 部分，添加：

| 变量名 | 值 | 说明 |
|--------|-----|------|
| `VITE_API_BASE_URL` | `https://your-backend-url.onrender.com/api` | 后端 API 地址（替换为你的 Render 后端 URL） |

### 3.4 部署

1. 点击 "Deploy"
2. 等待构建完成（约 2-3 分钟）
3. 部署成功后，你会得到一个 URL，例如：`https://book-manage.vercel.app`
4. **保存这个 URL**

### 3.5 更新 Vercel 配置（可选）

如果需要更精细的控制，可以编辑项目根目录的 `vercel.json` 文件，更新后端 URL：

```json
{
  "routes": [
    {
      "src": "/api/(.*)",
      "dest": "https://your-backend-url.onrender.com/api/$1"
    }
  ]
}
```

然后重新部署。

---

## ✅ 第四步：验证部署

### 4.1 测试前端

1. 访问你的 Vercel 前端地址
2. 尝试注册一个新用户
3. 尝试登录
4. 检查功能是否正常

### 4.2 检查后端日志

1. 在 Render Dashboard 中，点击你的后端服务
2. 查看 **Logs** 标签页
3. 确认没有错误信息

### 4.3 检查数据库

1. 在 Supabase Dashboard 中，点击 **Table Editor**
2. 查看 `user` 表，确认新注册的用户已保存

---

## 🔧 常见问题

### 问题 1：后端连接数据库失败

**解决方案**：
- 检查 `DATABASE_URL` 环境变量是否正确
- 确认 Supabase 数据库密码是否正确
- 检查 Supabase 项目的 **Settings** → **Database** → **Connection pooling** 是否启用

### 问题 2：前端无法访问后端 API

**解决方案**：
- 检查 `VITE_API_BASE_URL` 环境变量是否正确
- 确认后端 URL 是否可访问
- 检查 CORS 配置（后端已配置，但可以检查 Render 日志）

### 问题 3：Render 服务休眠

**免费计划**：Render 的免费服务在 15 分钟无活动后会休眠，首次访问需要等待约 30 秒唤醒。

**解决方案**：
- 升级到付费计划（$7/月）
- 或使用其他平台（如 Railway、Fly.io）

### 问题 4：数据库迁移问题

如果导入 SQL 脚本时出错：

1. 检查 SQL 语法是否正确
2. 确认表是否已存在（如果已存在，先删除再导入）
3. 在 Supabase SQL Editor 中逐步执行 SQL 语句

---

## 📝 环境变量参考

### 后端环境变量（Render）

```bash
DB_TYPE=postgres
DATABASE_URL=postgresql://postgres:password@db.xxx.supabase.co:5432/postgres
JWT_SECRET=your-jwt-secret-key
PORT=8080
ADMIN_EMAILS=admin@lib.com
SMTP_HOST=smtp.qq.com
SMTP_PORT=465
SMTP_USER=your-email@qq.com
SMTP_PASSWORD=your-email-password
```

### 前端环境变量（Vercel）

```bash
VITE_API_BASE_URL=https://your-backend-url.onrender.com/api
```

---

## 🔐 安全建议

1. **JWT_SECRET**：使用强随机字符串，不要使用默认值
2. **数据库密码**：使用强密码，定期更换
3. **邮箱密码**：如果使用 QQ 邮箱，建议使用授权码而非登录密码
4. **环境变量**：不要在代码中硬编码敏感信息
5. **HTTPS**：Vercel 和 Render 都自动提供 HTTPS

---

## 🎉 完成！

恭喜！你的图书管理系统已经成功部署到线上。

- **前端地址**: https://your-frontend.vercel.app
- **后端地址**: https://your-backend.onrender.com
- **数据库**: Supabase Dashboard

如有问题，请查看各平台的文档或联系支持。

---

## 📚 相关链接

- [Supabase 文档](https://supabase.com/docs)
- [Vercel 文档](https://vercel.com/docs)
- [Render 文档](https://render.com/docs)
- [项目 README](./README.md)


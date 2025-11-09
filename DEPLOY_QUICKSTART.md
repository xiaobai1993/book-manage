# 快速部署指南（5 分钟版）

## 🚀 超快速部署步骤

### 1️⃣ Supabase 数据库（2 分钟）

1. 访问 https://supabase.com 注册并创建项目
2. 在 SQL Editor 中执行 `data_postgresql.sql`
3. 在 Settings → Database 复制 `DATABASE_URL`

### 2️⃣ Render 后端（2 分钟）

1. 访问 https://render.com 注册并连接 GitHub
2. 创建 Web Service，选择你的仓库
3. 配置：
   - **Build**: `go mod download && go build -o book-manage`
   - **Start**: `./book-manage`
4. 添加环境变量：
   ```
   DB_TYPE=postgres
   DATABASE_URL=你的 Supabase 连接字符串
   JWT_SECRET=随机生成的密钥
   ```
5. 部署并复制后端 URL

### 3️⃣ Vercel 前端（1 分钟）

1. 访问 https://vercel.com 注册并连接 GitHub
2. 导入项目，设置 Root Directory 为 `frontend`
3. 添加环境变量：
   ```
   VITE_API_BASE_URL=https://你的后端URL.onrender.com/api
   ```
4. 部署完成！

---

## 📋 详细步骤请查看 [DEPLOY.md](./DEPLOY.md)


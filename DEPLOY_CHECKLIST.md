# 部署检查清单

在开始部署前，请确认以下事项：

## ✅ 代码准备

- [x] PostgreSQL 数据库迁移脚本已创建 (`data_postgresql.sql`)
- [x] Go 代码已支持 PostgreSQL (`database/db.go`)
- [x] 配置系统已支持环境变量 (`config/config.go`)
- [x] 前端已支持生产环境 API 地址 (`frontend/src/utils/request.js`)
- [x] 部署配置文件已创建 (`vercel.json`, `render.yaml`)

## 📝 部署前准备

### Supabase
- [ ] 注册 Supabase 账号
- [ ] 创建新项目
- [ ] 保存数据库密码
- [ ] 执行 `data_postgresql.sql` 脚本
- [ ] 复制 `DATABASE_URL` 连接字符串

### Render（后端）
- [ ] 注册 Render 账号
- [ ] 连接 GitHub 仓库
- [ ] 准备环境变量：
  - [ ] `DB_TYPE=postgres`
  - [ ] `DATABASE_URL`（从 Supabase 获取）
  - [ ] `JWT_SECRET`（生成强随机字符串）
  - [ ] `ADMIN_EMAILS`（管理员邮箱）
  - [ ] 其他可选配置（邮箱等）

### Vercel（前端）
- [ ] 注册 Vercel 账号
- [ ] 连接 GitHub 仓库
- [ ] 准备环境变量：
  - [ ] `VITE_API_BASE_URL`（后端 URL + `/api`）

## 🚀 部署步骤

### 第一步：Supabase 数据库
1. [ ] 创建 Supabase 项目
2. [ ] 执行 SQL 脚本
3. [ ] 获取连接字符串

### 第二步：Render 后端
1. [ ] 创建 Web Service
2. [ ] 配置构建和启动命令
3. [ ] 添加环境变量
4. [ ] 部署并获取 URL

### 第三步：Vercel 前端
1. [ ] 导入项目
2. [ ] 设置 Root Directory 为 `frontend`
3. [ ] 添加环境变量（包含后端 URL）
4. [ ] 部署

## ✅ 部署后验证

- [ ] 前端可以正常访问
- [ ] 可以注册新用户
- [ ] 可以登录
- [ ] 可以查看图书列表
- [ ] 管理员功能正常
- [ ] 数据库数据正确保存

## 🔧 故障排查

如果遇到问题，检查：

1. **后端无法连接数据库**
   - 检查 `DATABASE_URL` 是否正确
   - 检查 Supabase 项目是否正常运行
   - 查看 Render 日志

2. **前端无法访问后端**
   - 检查 `VITE_API_BASE_URL` 是否正确
   - 检查后端 URL 是否可访问
   - 检查 CORS 配置

3. **服务休眠（Render 免费版）**
   - 首次访问需要等待约 30 秒
   - 考虑升级到付费计划

---

详细步骤请参考 [DEPLOY.md](./DEPLOY.md)


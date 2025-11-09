# 本地 PostgreSQL 测试指南

## 📌 重要说明

**你不需要在本地切换到 PostgreSQL！**

- ✅ 本地开发可以继续使用 MySQL（默认配置）
- ✅ 部署到生产环境时会自动使用 PostgreSQL（通过环境变量）
- ✅ 代码已经同时支持两种数据库，无需修改

---

## 🤔 为什么不需要切换？

1. **代码已支持双数据库**：通过 `DB_TYPE` 环境变量自动切换
2. **本地开发环境**：继续使用 MySQL，配置简单，无需改动
3. **生产环境**：Render 会自动设置 `DB_TYPE=postgres`，使用 Supabase PostgreSQL

---

## 💡 如果你想本地测试 PostgreSQL（可选）

如果你想在本地测试 PostgreSQL 以确保兼容性，可以按以下步骤：

### 1. 安装 PostgreSQL

**macOS:**
```bash
brew install postgresql@14
brew services start postgresql@14
```

**或使用 Docker:**
```bash
docker run --name postgres-test -e POSTGRES_PASSWORD=123456 -p 5432:5432 -d postgres:14
```

### 2. 创建数据库

```bash
# 连接到 PostgreSQL
psql -U postgres

# 创建数据库
CREATE DATABASE library_management;

# 退出
\q
```

### 3. 导入数据

```bash
# 执行 PostgreSQL 迁移脚本
psql -U postgres -d library_management -f data_postgresql.sql
```

### 4. 配置环境变量

**方式一：使用环境变量（推荐）**

```bash
# 设置数据库类型
export DB_TYPE=postgres

# 设置数据库连接（如果使用配置文件）
export APP_ENV=env  # 使用 env.yaml，但需要修改为 PostgreSQL 配置

# 或者直接使用 DATABASE_URL（Supabase 格式）
export DATABASE_URL=postgresql://postgres:123456@localhost:5432/library_management
```

**方式二：修改配置文件**

1. 复制 `config/env_postgresql.yaml.example` 为 `config/env_postgresql.yaml`
2. 修改数据库连接信息
3. 设置环境变量：`export APP_ENV=env_postgresql`

### 5. 运行后端

```bash
# 确保设置了 DB_TYPE=postgres
export DB_TYPE=postgres
go run main.go
```

---

## ✅ 总结

**对于部署：**
- ✅ 不需要修改本地配置
- ✅ 不需要在本地安装 PostgreSQL
- ✅ 直接按照 `DEPLOY.md` 部署即可
- ✅ 生产环境会自动使用 PostgreSQL

**如果你想本地测试 PostgreSQL：**
- 按照上面的步骤操作
- 或者直接部署到生产环境测试

---

## 🔗 相关文档

- [部署指南](./DEPLOY.md) - 生产环境部署步骤
- [快速部署](./DEPLOY_QUICKSTART.md) - 5分钟快速部署


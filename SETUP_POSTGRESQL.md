# PostgreSQL 本地设置指南

## 📋 快速设置步骤

### 1. 安装 PostgreSQL

**macOS (使用 Homebrew):**
```bash
brew install postgresql@14
brew services start postgresql@14
```

**或者使用 Docker (推荐，更简单):**
```bash
docker run --name postgres-book-manage \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_USER=postgres \
  -p 5432:5432 \
  -d postgres:14
```

**Windows:**
- 下载并安装 [PostgreSQL](https://www.postgresql.org/download/windows/)
- 安装时设置密码为 `postgres`

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### 2. 创建数据库

**方式一：使用 psql 命令行**
```bash
# 连接到 PostgreSQL
psql -U postgres

# 创建数据库
CREATE DATABASE library_management;

# 退出
\q
```

**方式二：使用 Docker（如果使用 Docker）**
```bash
# 进入容器
docker exec -it postgres-book-manage psql -U postgres

# 创建数据库
CREATE DATABASE library_management;

# 退出
\q
```

### 3. 导入数据库结构

```bash
# 使用 psql 导入
psql -U postgres -d library_management -f data_postgresql.sql

# 或者使用 Docker
docker exec -i postgres-book-manage psql -U postgres -d library_management < data_postgresql.sql
```

### 4. 验证数据库

```bash
# 连接到数据库
psql -U postgres -d library_management

# 查看表
\dt

# 应该看到三个表：user, book, borrow_record
# 退出
\q
```

### 5. 运行项目

现在可以直接运行项目，代码已经配置为使用 PostgreSQL：

```bash
# 运行后端
go run main.go

# 或者使用 Makefile
make dev-backend
```

---

## 🔧 配置说明

项目已配置为默认使用 PostgreSQL，配置文件位于：
- `config/env.yaml` - 本地开发环境
- `config/dev.yaml` - 开发环境
- `config/prod.yaml` - 生产环境

所有配置文件都已设置为：
```yaml
database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"
  database: "library_management"
```

---

## 🐳 Docker 快速启动（推荐）

如果你使用 Docker，可以一键启动 PostgreSQL：

```bash
# 启动 PostgreSQL 容器
docker run --name postgres-book-manage \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_USER=postgres \
  -p 5432:5432 \
  -d postgres:14

# 等待几秒让数据库启动
sleep 5

# 创建数据库
docker exec -i postgres-book-manage psql -U postgres -c "CREATE DATABASE library_management;"

# 导入数据
docker exec -i postgres-book-manage psql -U postgres -d library_management < data_postgresql.sql

# 查看容器状态
docker ps | grep postgres-book-manage
```

停止和删除容器：
```bash
# 停止容器
docker stop postgres-book-manage

# 删除容器
docker rm postgres-book-manage
```

---

## ✅ 验证安装

运行以下命令验证数据库连接：

```bash
# 测试连接
psql -U postgres -d library_management -c "SELECT version();"

# 查看表
psql -U postgres -d library_management -c "\dt"

# 查看用户数据
psql -U postgres -d library_management -c "SELECT * FROM \"user\";"
```

---

## 🆘 常见问题

### 问题 1: 连接被拒绝

**错误**: `connection refused` 或 `could not connect to server`

**解决方案**:
- 确认 PostgreSQL 服务已启动
- macOS: `brew services list` 查看服务状态
- Linux: `sudo systemctl status postgresql`
- Docker: `docker ps` 查看容器是否运行

### 问题 2: 认证失败

**错误**: `password authentication failed`

**解决方案**:
- 确认密码是 `postgres`
- 如果修改了密码，需要更新配置文件中的密码

### 问题 3: 数据库不存在

**错误**: `database "library_management" does not exist`

**解决方案**:
```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE library_management;"
```

### 问题 4: 端口被占用

**错误**: `port 5432 is already in use`

**解决方案**:
- 检查是否有其他 PostgreSQL 实例在运行
- 或者修改配置文件中的端口号

---

## 📚 相关文档

- [PostgreSQL 官方文档](https://www.postgresql.org/docs/)
- [项目 README](./README.md)
- [部署指南](./DEPLOY.md)


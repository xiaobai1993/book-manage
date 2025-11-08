# 配置文件说明

## 环境配置

项目支持三个环境的配置文件：
- `env.yaml` - 开发环境（默认）
- `dev.yaml` - 开发环境
- `prod.yaml` - 生产环境

## 使用方法

### 1. 通过环境变量指定环境

```bash
# 使用env环境（默认）
export APP_ENV=env
go run main.go

# 使用dev环境
export APP_ENV=dev
go run main.go

# 使用prod环境
export APP_ENV=prod
go run main.go
```

### 2. 配置文件结构

```yaml
database:
  host: "localhost"
  port: "3306"
  user: "root"
  password: "123456"
  database: "library_management"

server:
  port: "8080"

jwt:
  secret: "book-manage-secret-key-2025"

email:
  smtp_host: ""
  smtp_port: ""
  smtp_user: ""
  smtp_password: ""

# 管理员邮箱白名单（优先判断）
admin_emails:
  - "824955445@qq.com"
  - "admin@lib.com"
```

## 管理员权限判断逻辑

1. **优先判断邮箱白名单**：如果用户邮箱在 `admin_emails` 列表中，直接认为是管理员
2. **其次检查数据库role字段**：如果不在白名单中，检查数据库中用户的 `role` 字段是否为 `admin`
3. **JWT Token中的role**：登录时会将确定的角色写入JWT token中

## 注意事项

1. 生产环境请修改JWT secret为强密码
2. 生产环境请修改数据库密码
3. 管理员邮箱白名单支持多个邮箱
4. 邮箱匹配不区分大小写

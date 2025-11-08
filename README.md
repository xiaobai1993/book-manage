# 图书管理系统后端服务

这是一个基于 Go 语言开发的图书管理系统后端服务，实现了完整的用户管理、图书管理和借阅管理功能。

## 功能特性

- ✅ 用户注册、登录、密码找回（邮箱验证码）
- ✅ JWT 认证和权限管理
- ✅ 图书信息的增删改查和搜索
- ✅ 图书借阅和归还
- ✅ 借阅记录查询（个人和全量）
- ✅ 管理员权限控制

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT
- **密码加密**: bcrypt

## 项目结构

```
book-manage/
├── main.go                 # 程序入口
├── config/                 # 配置模块
│   └── config.go
├── database/               # 数据库连接
│   └── db.go
├── models/                 # 数据模型
│   ├── user.go
│   ├── book.go
│   └── borrow_record.go
├── handlers/               # 请求处理器
│   ├── user.go
│   ├── book.go
│   └── borrow.go
├── middleware/             # 中间件
│   ├── auth.go
│   └── cors.go
├── services/               # 业务服务
│   └── email.go
└── utils/                  # 工具函数
    ├── jwt.go
    ├── response.go
    └── validator.go
```

## 安装和运行

### 1. 前置要求

- Go 1.21 或更高版本
- MySQL 5.7 或更高版本
- 已导入数据库脚本（data.sql）

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

编辑 `config/config.go` 文件，修改数据库连接信息（默认已配置为 localhost:3306，用户名 root，密码 123456）：

```go
Database: DatabaseConfig{
    Host:     "localhost",
    Port:     "3306",
    User:     "root",
    Password: "123456",
    Database: "library_management",
}
```

### 4. 运行服务

```bash
go run main.go
```

服务默认运行在 `http://localhost:8080`

## API 接口

所有接口均采用 POST 方法，请求头设置为 `Content-Type: application/json`。

详细的 API 文档请参考 `API.md` 文件。

### 主要接口列表

#### 用户管理
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录
- `POST /api/user/sendEmailCode` - 发送邮箱验证码
- `POST /api/user/forgetPassword` - 密码找回
- `POST /api/user/profile` - 获取个人信息（需登录）
- `POST /api/user/changePassword` - 修改密码（需登录）
- `POST /api/user/borrowRecords` - 获取个人借阅记录（需登录）

#### 图书管理
- `POST /api/book/search` - 图书搜索（需登录）
- `POST /api/book/detail` - 获取图书详情（需登录）
- `POST /api/book/add` - 添加图书（需管理员权限）
- `POST /api/book/edit` - 编辑图书（需管理员权限）
- `POST /api/book/delete` - 删除图书（需管理员权限）

#### 借阅管理
- `POST /api/borrow/borrow` - 借书（需登录）
- `POST /api/borrow/return` - 还书（需登录）
- `POST /api/borrow/records` - 获取个人借阅记录（需登录）
- `POST /api/borrow/allRecords` - 获取全量借阅记录（需管理员权限）

## 认证方式

支持三种方式传递 token：

1. **Header 方式（推荐）**:
   ```
   Authorization: Bearer <token>
   ```

2. **Query 参数方式**:
   ```
   ?token=<token>
   ```

3. **请求体方式**:
   ```json
   {
     "token": "<token>",
     ...
   }
   ```

## 测试账号

数据库初始化后，可以使用以下测试账号：

- **管理员账号**:
  - 邮箱: `admin@lib.com`
  - 密码: `12345678`

- **普通用户**:
  - 邮箱: `user1@lib.com`
  - 密码: `12345678`
  - 邮箱: `user2@lib.com`
  - 密码: `12345678`

> **注意**: 实际密码请参考数据库中的 bcrypt 哈希值。如果需要测试，请先使用注册功能创建新用户。

## 邮箱验证码

当前版本的邮箱验证码服务使用内存存储，验证码会打印到控制台。实际生产环境需要配置真实的邮件服务。

验证码规则：
- 6位数字
- 有效期30分钟
- 每分钟最多重发1次

## 错误码说明

| 错误码 | 含义 |
|--------|------|
| 0 | 成功 |
| 10001 | 参数错误 |
| 10002 | 邮箱格式错误 |
| 10003 | 邮箱已被注册 |
| 10004 | 验证码错误或已过期 |
| 10005 | 密码不一致 |
| 10006 | 密码长度不足 |
| 10007 | 邮箱或密码错误 |
| 10008 | 邮箱未注册 |
| 10009 | 权限不足 |
| 10010 | 图书不存在 |
| 10011 | 库存不足 |
| 10012 | 借阅已达上限 |
| 10013 | 该图书已存在未归还记录 |
| 10014 | 图书不可删除 |
| 10015 | 不存在此借阅记录 |
| 10016 | 图书已归还 |
| 10017 | ISBN已存在 |
| 10018 | 总数量不能小于已借出数量 |
| 10019 | 搜索关键词过短 |
| 10020 | 搜索无结果 |

## 开发说明

### 添加新接口

1. 在 `handlers/` 目录下创建或修改对应的 handler 文件
2. 在 `main.go` 中添加路由配置
3. 根据需要添加中间件（认证、权限等）

### 数据库迁移

当前版本直接使用 SQL 脚本初始化数据库。如需使用 GORM 的自动迁移功能，可以在 `database/db.go` 中添加：

```go
db.AutoMigrate(&models.User{}, &models.Book{}, &models.BorrowRecord{})
```

## 许可证

MIT License

## 联系方式

如有问题或建议，请联系项目维护者。

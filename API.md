# 图书管理系统 v1.0 API 接口文档

## 1. 文档信息
- **API版本**：v1.0
- **文档版本**：v1.0
- **创建日期**：2025年11月8日
- **最后更新**：2025年11月8日

## 2. 通用规范

### 2.1 请求规范
- 所有接口均采用 `POST` 方法
- 请求头设置：`Content-Type: application/json`
- 所有HTTP状态码均为 `200`，业务状态通过响应体中的 `code` 字段判断
- 请求参数统一以 JSON 格式传递

### 2.2 响应规范
所有响应体遵循以下格式：

```
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| 字段 | 类型 | 含义 | 
|------|------|------|
| code | int | 业务状态码，0表示成功，其他值表示错误 |
| message | string | 响应消息，成功为"success"，错误时为具体错误信息 |
| data | object | 响应数据，根据接口不同内容不同 |

### 2.3 公共结构

#### 用户信息结构（User）
| 字段 | 类型 | 含义 |
|------|------|------|
| id | int | 用户ID |
| email | string | 用户邮箱 |
| role | string | 用户角色（admin/user） |
| register_time | string | 注册时间 |
| status | string | 账户状态（normal/disabled） |

#### 图书信息结构（Book）
| 字段 | 类型 | 含义 |
|------|------|------|
| id | int | 图书ID |
| title | string | 书名 |
| author | string | 作者 |
| isbn | string | ISBN编号 |
| category | string | 图书分类 |
| total_quantity | int | 总数量 |
| available_quantity | int | 可借数量 |
| description | string | 图书描述 |
| create_time | string | 添加时间 |
| update_time | string | 更新时间 |

#### 借阅记录结构（Borrow Record）
| 字段 | 类型 | 含义 |
|------|------|------|
| id | int | 记录ID |
| user_id | int | 借阅用户ID |
| user_email | string | 借阅用户邮箱 |
| book_id | int | 借阅图书ID |
| book_title | string | 借阅图书名称 |
| borrow_date | string | 借阅日期 |
| due_date | string | 应还日期 |
| return_date | string | 实际归还日期（未归还可以为null） |
| status | string | 借阅状态（borrowed/returned） |

### 2.4 业务错误码定义

| 错误码 | 含义 | 备注 |
|--------|------|------|
| 0 | 成功 | 操作成功执行 |
| 10001 | 参数错误 | 请求参数缺失或格式错误 |
| 10002 | 邮箱格式错误 | 邮箱格式不符合标准格式 |
| 10003 | 邮箱已被注册 | 用户注册时邮箱已被注册 |
| 10004 | 验证码错误或已过期 | 验证码错误或已过期 |
| 10005 | 密码不一致 | 确认密码与原密码不一致 |
| 10006 | 密码长度不足 | 密码长度少于8位 |
| 10007 | 邮箱或密码错误 | 用户登录时邮箱或密码错误 |
| 10008 | 邮箱未注册 | 验证邮箱是否为注册用户 |
| 10009 | 权限不足 | 当前用户权限不足以执行操作 |
| 10010 | 图书不存在 | 操作的图书不存在 |
| 10011 | 库存不足 | 图书可借数量不足 |
| 10012 | 借阅已达上限 | 用户当前借阅量达到5本上限 |
| 10013 | 该图书已存在未归还记录 | 试图同时借阅已借阅的图书 |
| 10014 | 图书不可删除 | 图书存在未归还借阅记录，无法删除 |
| 10015 | 不存在此借阅记录 | 用户未借阅此图书 |
| 10016 | 图书已归还 | 图书已被归还，无需重复操作 |
| 10017 | ISBN已存在 | 添加图书时ISBN已存在 |
| 10018 | 总数量不能小于已借出数量 | 更新图书数量时不符合业务规则 |
| 10019 | 搜索关键词过短 | 关键词搜索时关键词长度不足 |
| 10020 | 搜索无结果 | 关键词搜索时未找到匹配结果 |

## 3. 用户管理模块

### 3.1 用户注册
- **接口地址**：`/api/user/register`
- **请求方法**：`POST`
- **权限校验**：无需登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| email | string | 是 | 用户邮箱 |
| password | string | 是 | 用户密码（至少8位） |
| confirm_password | string | 是 | 确认密码 |
| code | string | 是 | 邮箱验证码 |

#### 响应参数
返回公共响应格式，data字段为空对象 `{}`

#### 示例请求
```
POST /api/user/register
{
  "email": "user@example.com",
  "password": "password123",
  "confirm_password": "password123",
  "code": "123456"
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 3.2 用户登录
- **接口地址**：`/api/user/login`
- **请求方法**：`POST`
- **权限校验**：无需登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| email | string | 是 | 用户邮箱 |
| password | string | 是 | 用户密码 |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| user_info | object | 用户信息 |
| token | string | 登录凭证 |

#### 示例请求
```
POST /api/user/login
{
  "email": "admin@lib.com",
  "password": "password123"
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {
    "user_info": {
      "id": 1,
      "email": "admin@lib.com",
      "role": "admin",
      "register_time": "2025-11-08 10:00:00",
      "status": "normal"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 3.3 获取邮箱验证码
- **接口地址**：`/api/user/sendEmailCode`
- **请求方法**：`POST`
- **权限校验**：无需登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| email | string | 是 | 邮箱地址 |
| action | string | 是 | 操作类型（register/forget） |

#### 响应参数
返回公共响应格式，data字段为空对象 `{}`

#### 示例请求
```
POST /api/user/sendEmailCode
{
  "email": "user@example.com",
  "action": "register"
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 3.4 密码找回
- **接口地址**：`/api/user/forgetPassword`
- **请求方法**：`POST`
- **权限校验**：无需登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| email | string | 是 | 用户邮箱 |
| code | string | 是 | 邮箱验证码 |
| new_password | string | 是 | 新密码 |
| confirm_new_password | string | 是 | 确认新密码 |

#### 响应参数
返回公共响应格式，data字段为空对象 `{}`

#### 示例请求
```
POST /api/user/forgetPassword
{
  "email": "user@example.com",
  "code": "123456",
  "new_password": "newpassword123",
  "confirm_new_password": "newpassword123"
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 3.5 获取个人信息
- **接口地址**：`/api/user/profile`
- **请求方法**：`POST`
- **权限校验**：需要登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| user_info | object | 用户信息 |
| current_borrow_count | int | 当前借阅数量 |

#### 示例请求
```
POST /api/user/profile
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {
    "user_info": {
      "id": 2,
      "email": "user1@lib.com",
      "role": "user",
      "register_time": "2025-11-08 10:00:00",
      "status": "normal"
    },
    "current_borrow_count": 1
  }
}
```

### 3.6 修改密码
- **接口地址**：`/api/user/changePassword`
- **请求方法**：`POST`
- **权限校验**：需要登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| old_password | string | 是 | 原密码 |
| new_password | string | 是 | 新密码 |
| confirm_new_password | string | 是 | 确认新密码 |

#### 响应参数
返回公共响应格式，data字段为空对象 `{}`

#### 示例请求
```
POST /api/user/changePassword
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "old_password": "oldpassword123",
  "new_password": "newpassword123",
  "confirm_new_password": "newpassword123"
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 3.7 获取个人借阅记录
- **接口地址**：`/api/user/borrowRecords`
- **请求方法**：`POST`
- **权限校验**：需要登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| status | string | 否 | 筛选状态（borrowed/returned/all） |
| page | int | 否 | 页码 |
| limit | int | 否 | 每页数量 |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| records | array | 借阅记录列表 |
| total | int | 总记录数 |
| page | int | 当前页码 |
| limit | int | 每页数量 |

#### 示例请求
```
POST /api/user/borrowRecords
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "status": "borrowed",
  "page": 1,
  "limit": 10
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {
    "records": [
      {
        "id": 1,
        "book_id": 4,
        "book_title": "小王子",
        "borrow_date": "2025-10-01 10:30:00",
        "due_date": "2025-10-31 10:30:00",
        "return_date": null,
        "status": "borrowed"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

## 4. 图书管理模块
- **权限说明**：图书信息的添加、编辑、删除功能仅管理员可见并操作；普通用户仅可查看与检索图书信息。

### 4.1 添加图书
- **接口地址**：`/api/book/add`
- **请求方法**：`POST`
- **权限校验**：需要管理员权限

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| title | string | 是 | 书名 |
| author | string | 是 | 作者 |
| isbn | string | 是 | ISBN编号 |
| category | string | 是 | 图书分类 |
| total_quantity | int | 是 | 总数量 |
| description | string | 否 | 图书描述 |

#### 响应参数
返回公共响应格式，data字段为空对象 `{}`

#### 示例请求
```
POST /api/book/add
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "title": "三体",
  "author": "刘慈欣",
  "isbn": "9787536692930",
  "category": "科幻",
  "total_quantity": 5,
  "description": "地球文明向宇宙发出了神秘信号..."
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 4.2 编辑图书
- **接口地址**：`/api/book/edit`
- **请求方法**：`POST`
- **权限校验**：需要管理员权限

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| id | int | 是 | 图书ID |
| title | string | 否 | 书名 |
| author | string | 否 | 作者 |
| isbn | string | 否 | ISBN编号 |
| category | string | 否 | 图书分类 |
| total_quantity | int | 否 | 总数量 |
| description | string | 否 | 图书描述 |

#### 响应参数
返回公共响应格式，data字段为空对象 `{}`

#### 示例请求
```
POST /api/book/edit
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "id": 1,
  "title": "三体",
  "author": "刘慈欣",
  "isbn": "9787536692930",
  "category": "科幻",
  "total_quantity": 6,
  "description": "地球文明向宇宙发出了神秘信号..."
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 4.3 删除图书
- **接口地址**：`/api/book/delete`
- **请求方法**：`POST`
- **权限校验**：需要管理员权限

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| id | int | 是 | 图书ID |

#### 响应参数
返回公共响应格式，data字段为空对象 `{}`

#### 示例请求
```
POST /api/book/delete
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "id": 1
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 4.4 获取图书详情
- **接口地址**：`/api/book/detail`
- **请求方法**：`POST`
- **权限校验**：需要登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| id | int | 是 | 图书ID |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| book | object | 图书信息 |

#### 示例请求
```
POST /api/book/detail
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "id": 1
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {
    "book": {
      "id": 1,
      "title": "三体",
      "author": "刘慈欣",
      "isbn": "9787536692930",
      "category": "科幻",
      "total_quantity": 5,
      "available_quantity": 5,
      "description": "地球文明向宇宙发出了神秘信号...",
      "create_time": "2025-11-08 10:00:00",
      "update_time": "2025-11-08 10:00:00"
    }
  }
}
```

### 4.5 图书搜索
- **接口地址**：`/api/book/search`
- **请求方法**：`POST`
- **权限校验**：需要登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| keyword | string | 否 | 搜索关键词（书名或作者） |
| category | string | 否 | 图书分类筛选 |
| page | int | 否 | 页码 |
| limit | int | 否 | 每页数量 |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| books | array | 图书列表 |
| total | int | 总记录数 |
| page | int | 当前页码 |
| limit | int | 每页数量 |

#### 示例请求
```
POST /api/book/search
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "keyword": "三体",
  "page": 1,
  "limit": 10
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {
    "books": [
      {
        "id": 1,
        "title": "三体",
        "author": "刘慈欣",
        "isbn": "9787536692930",
        "category": "科幻",
        "total_quantity": 5,
        "available_quantity": 5,
        "description": "地球文明向宇宙发出了神秘信号...",
        "create_time": "2025-11-08 10:00:00"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

## 5. 借阅管理模块

### 5.1 借书
- **接口地址**：`/api/borrow/borrow`
- **请求方法**：`POST`
- **权限校验**：需要登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| book_id | int | 是 | 图书ID |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| borrow_date | string | 借阅日期 |
| due_date | string | 应还日期 |

#### 示例请求
```
POST /api/borrow/borrow
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "book_id": 1
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {
    "borrow_date": "2025-11-08 14:30:00",
    "due_date": "2025-12-08 14:30:00"
  }
}
```

### 5.2 还书
- **接口地址**：`/api/borrow/return`
- **请求方法**：`POST`
- **权限校验**：需要登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| book_id | int | 是 | 图书ID |

#### 响应参数
返回公共响应格式，data字段为空对象 `{}`

#### 示例请求
```
POST /api/borrow/return
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "book_id": 4
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 5.3 获取借阅记录（普通用户）
- **接口地址**：`/api/borrow/records`
- **请求方法**：`POST`
- **权限校验**：需要登录

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| status | string | 否 | 筛选状态（borrowed/returned/all） |
| page | int | 否 | 页码 |
| limit | int | 否 | 每页数量 |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| records | array | 借阅记录列表 |
| total | int | 总记录数 |
| page | int | 当前页码 |
| limit | int | 每页数量 |

#### 示例请求
```
POST /api/borrow/records
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "status": "borrowed",
  "page": 1,
  "limit": 10
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {
    "records": [
      {
        "id": 1,
        "book_id": 4,
        "book_title": "小王子",
        "borrow_date": "2025-10-01 10:30:00",
        "due_date": "2025-10-31 10:30:00",
        "return_date": null,
        "status": "borrowed"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

### 5.4 获取全量借阅记录（仅管理员）
- **接口地址**：`/api/borrow/allRecords`
- **请求方法**：`POST`
- **权限校验**：需要管理员权限

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| user_email | string | 否 | 用户邮箱筛选 |
| book_title | string | 否 | 图书名称筛选 |
| status | string | 否 | 筛选状态（borrowed/returned/all） |
| page | int | 否 | 页码 |
| limit | int | 否 | 每页数量 |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| records | array | 借阅记录列表 |
| total | int | 总记录数 |
| page | int | 当前页码 |
| limit | int | 每页数量 |

#### 示例请求
```
POST /api/borrow/allRecords
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "status": "borrowed",
  "page": 1,
  "limit": 10
}
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {
    "records": [
      {
        "id": 1,
        "user_id": 2,
        "user_email": "user1@lib.com",
        "book_id": 4,
        "book_title": "小王子",
        "borrow_date": "2025-10-01 10:30:00",
        "due_date": "2025-10-31 10:30:00",
        "return_date": null,
        "status": "borrowed"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

## 6. 业务规则与约束

### 6.1 借阅规则
- 普通用户单次最多借阅5本图书
- 每本图书借阅期限为30天（当前版本仅记录应还日期，无逾期罚款机制）
- 暂不支持图书预约功能

### 6.2 权限规则
| 功能模块               | 普通用户 | 管理员 |
| ---------------------- | -------- | ------ |
| 注册/登录/密码找回     | √        | √      |
| 个人信息管理           | √        | √      |
| 图书检索与查看         | √        | √      |
| 图书添加/编辑/删除     | ×        | √      |
| 借还书操作             | √        | √      |
| 个人借阅记录查询       | √        | √      |
| 全量借阅记录查询       | ×        | √      |

### 6.3 数据约束
- 用户邮箱：需符合标准邮箱格式，且在系统内唯一
- 用户密码：长度≥8位，存储时需通过bcrypt算法加密
- 图书ISBN：系统内唯一，需符合ISBN编码规则
- 图书数量：总数量、可借数量均为非负整数

## 7. 完整API接口列表

| 接口名称 | 接口地址 | 权限要求 | 功能模块 |
|----------|----------|----------|----------|
| 用户注册 | `/api/user/register` | 无需登录 | 用户管理 |
| 用户登录 | `/api/user/login` | 无需登录 | 用户管理 |
| 发送邮箱验证码 | `/api/user/sendEmailCode` | 无需登录 | 用户管理 |
| 密码找回 | `/api/user/forgetPassword` | 无需登录 | 用户管理 |
| 获取个人信息 | `/api/user/profile` | 需要登录 | 用户管理 |
| 修改密码 | `/api/user/changePassword` | 需要登录 | 用户管理 |
| 获取个人借阅记录 | `/api/user/borrowRecords` | 需要登录 | 用户管理 |
| 添加图书 | `/api/book/add` | 需要管理员权限 | 图书管理 |
| 编辑图书 | `/api/book/edit` | 需要管理员权限 | 图书管理 |
| 删除图书 | `/api/book/delete` | 需要管理员权限 | 图书管理 |
| 获取图书详情 | `/api/book/detail` | 需要登录 | 图书管理 |
| 图书搜索 | `/api/book/search` | 需要登录 | 图书管理 |
| 借书 | `/api/borrow/borrow` | 需要登录 | 借阅管理 |
| 还书 | `/api/borrow/return` | 需要登录 | 借阅管理 |
| 获取个人借阅记录 | `/api/borrow/records` | 需要登录 | 借阅管理 |
| 获取全量借阅记录 | `/api/borrow/allRecords` | 需要管理员权限 | 借阅管理 |
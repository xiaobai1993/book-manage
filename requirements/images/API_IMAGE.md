# 图书图片管理 API 接口文档

## 1. 文档信息
- **API版本**：v1.0
- **文档版本**：v1.0
- **创建日期**：2025年1月
- **关联文档**：`../../API.md`（主接口文档）

## 2. 通用规范

### 2.1 请求规范
- 图片上传接口采用 `POST` 方法，使用 `multipart/form-data` 格式
- 其他接口采用 `POST` 方法，使用 `application/json` 格式
- 请求头设置：根据接口类型设置相应的 `Content-Type`
- 所有HTTP状态码均为 `200`，业务状态通过响应体中的 `code` 字段判断

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

### 2.3 业务错误码定义

| 错误码 | 含义 | 备注 |
|--------|------|------|
| 0 | 成功 | 操作成功执行 |
| 10001 | 参数错误 | 请求参数缺失或格式错误 |
| 10010 | 图书不存在 | 操作的图书不存在 |
| 10021 | 图片格式不支持 | 上传的图片格式不在允许范围内（仅支持 JPG、PNG、WebP、GIF） |
| 10022 | 图片大小超限 | 上传的图片大小超过限制（5MB） |
| 10023 | 图片上传失败 | 图片上传到存储服务失败 |
| 10024 | 图片不存在 | 要删除的图片不存在 |
| 10009 | 权限不足 | 当前用户权限不足以执行操作（需要管理员权限） |

## 3. 图片管理接口

### 3.1 上传图书封面
- **接口地址**：`/api/book/uploadCover`
- **请求方法**：`POST`
- **权限校验**：需要管理员权限
- **Content-Type**：`multipart/form-data`

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| book_id | int | 是 | 图书ID |
| image | file | 是 | 图片文件（JPG、PNG、WebP、GIF，最大 5MB） |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| image_url | string | 上传后的图片URL |
| book_id | int | 图书ID |

#### 示例请求
```
POST /api/book/uploadCover
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="token"

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="book_id"

1
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="image"; filename="cover.jpg"
Content-Type: image/jpeg

[图片二进制数据]
------WebKitFormBoundary7MA4YWxkTrZu0gW--
```

#### 示例响应
```
{
  "code": 0,
  "message": "success",
  "data": {
    "image_url": "https://your-r2-bucket.r2.cloudflarestorage.com/book-covers/1_a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg",
    "book_id": 1
  }
}
```

#### 错误响应示例
```
{
  "code": 10021,
  "message": "图片格式不支持，仅支持 JPG、PNG、WebP、GIF 格式",
  "data": {}
}
```

```
{
  "code": 10022,
  "message": "图片大小不能超过 5MB",
  "data": {}
}
```

### 3.2 删除图书封面
- **接口地址**：`/api/book/deleteCover`
- **请求方法**：`POST`
- **权限校验**：需要管理员权限

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |
| book_id | int | 是 | 图书ID |

#### 响应参数
返回公共响应格式，data字段为空对象 `{}`

#### 示例请求
```
POST /api/book/deleteCover
Content-Type: application/json

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
  "data": {}
}
```

#### 错误响应示例
```
{
  "code": 10024,
  "message": "该图书没有封面图片",
  "data": {}
}
```

### 3.3 获取图片上传配置（可选）
- **接口地址**：`/api/book/uploadConfig`
- **请求方法**：`POST`
- **权限校验**：需要管理员权限
- **功能说明**：获取图片上传的配置信息，如允许的格式、大小限制等（用于前端验证）

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| token | string | 是 | 登录凭证 |

#### 响应参数
| 参数名 | 类型 | 说明 |
|--------|------|------|
| allowed_formats | array | 允许的图片格式列表 |
| max_size | int | 最大文件大小（字节） |
| max_size_mb | float | 最大文件大小（MB） |

#### 示例请求
```
POST /api/book/uploadConfig
Content-Type: application/json

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
    "allowed_formats": ["jpg", "jpeg", "png", "webp", "gif"],
    "max_size": 5242880,
    "max_size_mb": 5.0
  }
}
```

## 4. 图片访问说明

### 4.1 图片URL格式
图片上传成功后，系统会返回一个公开访问的 URL，格式如下：
```
https://{account-id}.r2.cloudflarestorage.com/{bucket-name}/book-covers/{book_id}_{uuid}.{ext}
```

### 4.2 图片访问方式
- 图片通过 Cloudflare R2 的公开 URL 直接访问
- 支持 HTTPS 访问
- 支持 CDN 加速（如果配置了 Cloudflare CDN）

### 4.3 默认占位图
如果图书没有封面图片，前端应显示默认占位图。建议使用以下占位图 URL 或本地资源：
- 本地占位图：`/assets/images/book-placeholder.png`
- 或使用在线占位图服务

## 5. 图片格式与大小限制

### 5.1 支持的图片格式
- JPEG / JPG
- PNG
- WebP
- GIF

### 5.2 大小限制
- 最大文件大小：5MB（5,242,880 字节）
- 建议图片尺寸：800x1200 像素（竖版封面）

### 5.3 图片处理建议
- 前端上传前可进行图片压缩
- 后端可选择性实现图片自动压缩（非必需）
- 建议使用 WebP 格式以减小文件大小

## 6. 完整接口列表

| 接口名称 | 接口地址 | 权限要求 | 功能说明 |
|----------|----------|----------|----------|
| 上传图书封面 | `/api/book/uploadCover` | 需要管理员权限 | 上传图书封面图片到 Cloudflare R2 |
| 删除图书封面 | `/api/book/deleteCover` | 需要管理员权限 | 删除图书封面图片 |
| 获取上传配置 | `/api/book/uploadConfig` | 需要管理员权限 | 获取图片上传配置信息（可选） |

## 7. 注意事项

### 7.1 权限要求
- 所有图片管理接口都需要管理员权限
- 普通用户只能查看图片，不能上传或删除

### 7.2 文件命名规则
- 系统自动生成唯一文件名：`{book_id}_{uuid}.{ext}`
- UUID 确保文件名唯一性
- 保留原始文件扩展名

### 7.3 旧图片处理
- 更新图片时，系统会自动删除旧图片
- 删除图书封面时，会同时删除 R2 中的图片文件和数据库记录

### 7.4 错误处理
- 上传失败时，不会创建数据库记录
- 删除失败时，会记录错误日志，但不会影响数据库操作
- 建议实现重试机制

## 8. 测试建议

### 8.1 功能测试
- 测试上传各种格式的图片
- 测试上传超大文件（应返回错误）
- 测试上传不支持格式（应返回错误）
- 测试删除不存在的图片（应返回错误）
- 测试权限验证（普通用户应无法上传）

### 8.2 性能测试
- 测试并发上传（最多 10 个并发）
- 测试大文件上传（接近 5MB）
- 测试图片访问速度

### 8.3 边界测试
- 测试空文件上传
- 测试损坏的图片文件
- 测试不存在的图书ID


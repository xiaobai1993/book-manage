# Cloudflare R2 图片存储实施方案

## 1. 方案概述

### 1.1 为什么选择 Cloudflare R2
- **免费额度**：每月 10GB 存储空间，1000 万次读取操作，100 万次写入操作
- **无出口费用**：与 Cloudflare CDN 集成，无数据出口费用
- **S3 兼容**：完全兼容 AWS S3 API，易于集成
- **全球 CDN**：自动通过 Cloudflare 全球网络加速
- **高可用性**：99.9% 可用性保证

### 1.2 方案架构
```
前端上传图片 
    ↓
后端 API 接收
    ↓
验证图片格式和大小
    ↓
上传到 Cloudflare R2
    ↓
获取公开 URL
    ↓
保存 URL 到数据库
    ↓
返回 URL 给前端
```

## 2. Cloudflare R2 配置步骤

### 2.1 创建 Cloudflare 账户
1. 访问 https://dash.cloudflare.com/sign-up
2. 注册或登录 Cloudflare 账户（免费账户即可）

### 2.2 创建 R2 存储桶
1. 登录 Cloudflare Dashboard
2. 进入 **R2** 服务（左侧菜单）
3. 点击 **Create bucket**（创建存储桶）
4. 输入存储桶名称：`book-covers-prod`（生产环境）或 `book-covers-dev`（开发环境）
5. 选择位置：选择离你最近的区域（如 `apac` 亚太地区）
6. 点击 **Create bucket** 完成创建

### 2.3 创建 API Token
1. 在 R2 页面，点击右上角 **Manage R2 API Tokens**
2. 点击 **Create API token**
3. 配置 Token：
   - **Token name**：`book-manage-r2-token`
   - **Permissions**：选择 **Object Read & Write**
   - **TTL**：选择 **Never expire**（或设置合适的过期时间）
   - **R2 Token Scopes**：选择你创建的存储桶
4. 点击 **Create API Token**
5. **重要**：复制并保存以下信息（只显示一次）：
   - **Access Key ID**
   - **Secret Access Key**

### 2.4 配置公开访问（可选）
如果需要通过公开 URL 直接访问图片：

1. 在存储桶页面，点击 **Settings**（设置）
2. 找到 **Public Access** 部分
3. 可以选择：
   - **Custom Domain**：使用自定义域名（需要配置 DNS）
   - **R2.dev Subdomain**：使用 Cloudflare 提供的子域名（最简单）

**注意**：如果使用 R2.dev Subdomain，需要：
1. 在存储桶设置中启用 **R2.dev Subdomain**
2. 系统会自动分配一个子域名，如：`https://pub-xxxxx.r2.dev`
3. 图片 URL 格式：`https://pub-xxxxx.r2.dev/book-covers/{filename}`

## 3. 后端实现方案

### 3.1 安装依赖
使用 Go 语言的 AWS SDK（兼容 S3 API）：

```bash
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/s3
go get github.com/google/uuid
```

### 3.2 配置文件更新
在 `config/config.go` 中添加 R2 配置：

```go
// CloudflareR2Config R2配置
type CloudflareR2Config struct {
    AccountID       string `yaml:"account_id"`
    AccessKeyID     string `yaml:"access_key_id"`
    SecretAccessKey string `yaml:"secret_access_key"`
    BucketName      string `yaml:"bucket_name"`
    PublicURL       string `yaml:"public_url"` // R2.dev 子域名或自定义域名
    Region          string `yaml:"region"`     // 默认: auto
}
```

在 `config/env.yaml` 中添加配置：

```yaml
cloudflare_r2:
  account_id: "your-account-id"  # 从 Cloudflare Dashboard 获取
  access_key_id: "your-access-key-id"  # 从 API Token 获取
  secret_access_key: "your-secret-access-key"  # 从 API Token 获取
  bucket_name: "book-covers-dev"
  public_url: "https://pub-xxxxx.r2.dev"  # R2.dev 子域名
  region: "auto"
```

### 3.3 环境变量配置（生产环境）
在生产环境中，通过环境变量配置（更安全）：

```bash
export R2_ACCOUNT_ID="your-account-id"
export R2_ACCESS_KEY_ID="your-access-key-id"
export R2_SECRET_ACCESS_KEY="your-secret-access-key"
export R2_BUCKET_NAME="book-covers-prod"
export R2_PUBLIC_URL="https://pub-xxxxx.r2.dev"
export R2_REGION="auto"
```

### 3.4 创建 R2 服务模块
创建 `services/r2.go` 文件：

```go
package services

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "path/filepath"
    "strings"
    
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/google/uuid"
)

type R2Service struct {
    client   *s3.Client
    bucket   string
    publicURL string
}

var r2Service *R2Service

func InitR2Service(cfg *config.CloudflareR2Config) error {
    // 创建 S3 兼容的配置
    r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
        return aws.Endpoint{
            URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID),
        }, nil
    })

    cfgOptions := []func(*config.LoadOptions) error{
        config.WithEndpointResolverWithOptions(r2Resolver),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            cfg.AccessKeyID,
            cfg.SecretAccessKey,
            "",
        )),
        config.WithRegion(cfg.Region),
    }

    awsCfg, err := config.LoadDefaultConfig(context.TODO(), cfgOptions...)
    if err != nil {
        return fmt.Errorf("failed to load AWS config: %w", err)
    }

    r2Service = &R2Service{
        client:    s3.NewFromConfig(awsCfg),
        bucket:    cfg.BucketName,
        publicURL: cfg.PublicURL,
    }

    return nil
}

func GetR2Service() *R2Service {
    return r2Service
}

// UploadImage 上传图片到 R2
func (s *R2Service) UploadImage(bookID int, imageData []byte, filename string) (string, error) {
    // 生成唯一文件名
    ext := strings.ToLower(filepath.Ext(filename))
    uniqueFilename := fmt.Sprintf("%d_%s%s", bookID, uuid.New().String(), ext)
    key := fmt.Sprintf("book-covers/%s", uniqueFilename)

    // 上传到 R2
    _, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        bytes.NewReader(imageData),
        ContentType: aws.String(getContentType(ext)),
    })
    if err != nil {
        return "", fmt.Errorf("failed to upload image: %w", err)
    }

    // 返回公开 URL
    imageURL := fmt.Sprintf("%s/%s", s.publicURL, key)
    return imageURL, nil
}

// DeleteImage 从 R2 删除图片
func (s *R2Service) DeleteImage(imageURL string) error {
    // 从 URL 中提取 key
    key := strings.TrimPrefix(imageURL, s.publicURL+"/")
    
    _, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return fmt.Errorf("failed to delete image: %w", err)
    }

    return nil
}

// getContentType 根据文件扩展名获取 Content-Type
func getContentType(ext string) string {
    ext = strings.ToLower(ext)
    switch ext {
    case ".jpg", ".jpeg":
        return "image/jpeg"
    case ".png":
        return "image/png"
    case ".webp":
        return "image/webp"
    case ".gif":
        return "image/gif"
    default:
        return "application/octet-stream"
    }
}
```

### 3.5 数据库迁移
在 `book` 表中添加 `cover_image_url` 字段：

```sql
ALTER TABLE book ADD COLUMN cover_image_url VARCHAR(500) DEFAULT NULL;
```

或使用 GORM 自动迁移（在 `models/book.go` 中添加字段）：

```go
CoverImageURL string `gorm:"type:varchar(500);default:null" json:"cover_image_url"`
```

### 3.6 创建图片上传 Handler
在 `handlers/book.go` 中添加：

```go
// UploadCoverRequest 上传封面请求
type UploadCoverRequest struct {
    Token  string `form:"token"`
    BookID int    `form:"book_id" binding:"required"`
}

// UploadCover 上传图书封面
func UploadCover(c *gin.Context) {
    var req UploadCoverRequest
    if err := c.ShouldBind(&req); err != nil {
        utils.Error(c, 10001, "参数错误")
        return
    }

    // 验证管理员权限（通过中间件）
    // ...

    // 获取上传的文件
    file, err := c.FormFile("image")
    if err != nil {
        utils.Error(c, 10001, "请选择图片文件")
        return
    }

    // 验证文件格式
    ext := strings.ToLower(filepath.Ext(file.Filename))
    allowedExts := []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
    if !contains(allowedExts, ext) {
        utils.Error(c, 10021, "图片格式不支持，仅支持 JPG、PNG、WebP、GIF 格式")
        return
    }

    // 验证文件大小（5MB）
    if file.Size > 5*1024*1024 {
        utils.Error(c, 10022, "图片大小不能超过 5MB")
        return
    }

    // 打开文件
    src, err := file.Open()
    if err != nil {
        utils.Error(c, 10001, "无法读取图片文件")
        return
    }
    defer src.Close()

    // 读取文件内容
    imageData, err := io.ReadAll(src)
    if err != nil {
        utils.Error(c, 10001, "无法读取图片文件")
        return
    }

    // 检查图书是否存在
    db := database.GetDB()
    var book models.Book
    if err := db.First(&book, req.BookID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.Error(c, 10010, "图书不存在")
        } else {
            utils.Error(c, 10001, "查询图书失败")
        }
        return
    }

    // 如果已有图片，先删除旧图片
    if book.CoverImageURL != "" {
        r2Service := services.GetR2Service()
        if r2Service != nil {
            _ = r2Service.DeleteImage(book.CoverImageURL) // 忽略删除错误
        }
    }

    // 上传到 R2
    r2Service := services.GetR2Service()
    if r2Service == nil {
        utils.Error(c, 10023, "图片存储服务未配置")
        return
    }

    imageURL, err := r2Service.UploadImage(req.BookID, imageData, file.Filename)
    if err != nil {
        utils.Error(c, 10023, "图片上传失败")
        return
    }

    // 更新数据库
    book.CoverImageURL = imageURL
    book.UpdateTime = time.Now()
    if err := db.Save(&book).Error; err != nil {
        utils.Error(c, 10001, "更新图书记录失败")
        return
    }

    utils.Success(c, map[string]interface{}{
        "image_url": imageURL,
        "book_id":   req.BookID,
    })
}

// DeleteCoverRequest 删除封面请求
type DeleteCoverRequest struct {
    Token  string `json:"token"`
    BookID int    `json:"book_id" binding:"required"`
}

// DeleteCover 删除图书封面
func DeleteCover(c *gin.Context) {
    var req DeleteCoverRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.Error(c, 10001, "参数错误")
        return
    }

    // 验证管理员权限（通过中间件）
    // ...

    db := database.GetDB()

    // 查找图书
    var book models.Book
    if err := db.First(&book, req.BookID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.Error(c, 10010, "图书不存在")
        } else {
            utils.Error(c, 10001, "查询图书失败")
        }
        return
    }

    // 检查是否有图片
    if book.CoverImageURL == "" {
        utils.Error(c, 10024, "该图书没有封面图片")
        return
    }

    // 从 R2 删除图片
    r2Service := services.GetR2Service()
    if r2Service != nil {
        if err := r2Service.DeleteImage(book.CoverImageURL); err != nil {
            // 记录错误但不阻止删除数据库记录
            // log.Printf("Failed to delete image from R2: %v", err)
        }
    }

    // 清空数据库记录
    book.CoverImageURL = ""
    book.UpdateTime = time.Now()
    if err := db.Save(&book).Error; err != nil {
        utils.Error(c, 10001, "更新图书记录失败")
        return
    }

    utils.Success(c, map[string]interface{}{})
}

// contains 检查字符串是否在切片中
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

### 3.7 添加路由
在 `main.go` 中添加路由：

```go
// 图片上传（需要管理员权限）
bookGroup.POST("/uploadCover", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UploadCover)
bookGroup.POST("/deleteCover", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteCover)
```

## 4. 前端实现方案

### 4.1 图片上传组件
在 `frontend/src/components/` 创建 `ImageUpload.vue`：

```vue
<template>
  <div class="image-upload">
    <input
      type="file"
      ref="fileInput"
      accept="image/jpeg,image/png,image/webp,image/gif"
      @change="handleFileChange"
      style="display: none"
    />
    <div
      class="upload-area"
      @click="$refs.fileInput.click()"
      @dragover.prevent
      @drop.prevent="handleDrop"
    >
      <img v-if="previewUrl" :src="previewUrl" alt="预览" />
      <div v-else class="upload-placeholder">
        <p>点击或拖拽上传图片</p>
        <p class="hint">支持 JPG、PNG、WebP、GIF，最大 5MB</p>
      </div>
    </div>
    <button @click="upload" :disabled="!file || uploading">
      {{ uploading ? '上传中...' : '上传' }}
    </button>
  </div>
</template>

<script>
import { uploadCover } from '@/api/book'

export default {
  props: {
    bookId: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      file: null,
      previewUrl: null,
      uploading: false
    }
  },
  methods: {
    handleFileChange(e) {
      const file = e.target.files[0]
      if (file) {
        this.validateAndSetFile(file)
      }
    },
    handleDrop(e) {
      const file = e.dataTransfer.files[0]
      if (file) {
        this.validateAndSetFile(file)
      }
    },
    validateAndSetFile(file) {
      // 验证格式
      const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
      if (!allowedTypes.includes(file.type)) {
        this.$message.error('仅支持 JPG、PNG、WebP、GIF 格式')
        return
      }
      
      // 验证大小
      if (file.size > 5 * 1024 * 1024) {
        this.$message.error('图片大小不能超过 5MB')
        return
      }
      
      this.file = file
      
      // 预览
      const reader = new FileReader()
      reader.onload = (e) => {
        this.previewUrl = e.target.result
      }
      reader.readAsDataURL(file)
    },
    async upload() {
      if (!this.file) return
      
      this.uploading = true
      const formData = new FormData()
      formData.append('token', this.$store.state.user.token)
      formData.append('book_id', this.bookId)
      formData.append('image', this.file)
      
      try {
        const res = await uploadCover(formData)
        if (res.code === 0) {
          this.$message.success('上传成功')
          this.$emit('uploaded', res.data.image_url)
        }
      } catch (error) {
        this.$message.error(error.message || '上传失败')
      } finally {
        this.uploading = false
      }
    }
  }
}
</script>
```

### 4.2 API 调用
在 `frontend/src/api/book.js` 中添加：

```javascript
export function uploadCover(formData) {
  return request({
    url: '/api/book/uploadCover',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function deleteCover(data) {
  return request({
    url: '/api/book/deleteCover',
    method: 'post',
    data
  })
}
```

## 5. 部署注意事项

### 5.1 环境变量配置
在生产环境（如 Render、Vercel 等）中，通过环境变量配置 R2 凭证：

```bash
R2_ACCOUNT_ID=your-account-id
R2_ACCESS_KEY_ID=your-access-key-id
R2_SECRET_ACCESS_KEY=your-secret-access-key
R2_BUCKET_NAME=book-covers-prod
R2_PUBLIC_URL=https://pub-xxxxx.r2.dev
R2_REGION=auto
```

### 5.2 安全建议
1. **不要将 R2 凭证提交到代码仓库**
2. **使用环境变量存储敏感信息**
3. **定期轮换 API Token**
4. **限制 API Token 权限范围**

### 5.3 监控与日志
- 监控 R2 使用量（存储空间、请求次数）
- 记录图片上传/删除操作日志
- 监控图片访问错误率

## 6. 成本估算

### 6.1 免费额度
- **存储空间**：10GB/月（免费）
- **读取操作**：1000 万次/月（免费）
- **写入操作**：100 万次/月（免费）

### 6.2 典型使用场景
假设：
- 每张图片平均 500KB
- 1000 本图书，每本 1 张封面
- 总存储：约 500MB（远低于 10GB 免费额度）
- 每月图片访问：10 万次（远低于 1000 万次免费额度）

**结论**：对于中小型图书管理系统，完全在免费额度内。

## 7. 故障处理

### 7.1 R2 服务不可用
- 实现重试机制
- 记录错误日志
- 返回友好的错误提示

### 7.2 图片上传失败
- 验证网络连接
- 检查 R2 凭证配置
- 检查存储桶权限

### 7.3 图片访问失败
- 检查公开访问配置
- 检查 URL 格式
- 实现默认占位图

## 8. 后续优化

### 8.1 图片压缩
- 上传时自动压缩
- 生成多种尺寸（缩略图、中等尺寸）
- 使用 WebP 格式

### 8.2 CDN 加速
- 配置 Cloudflare CDN
- 使用自定义域名
- 启用缓存策略

### 8.3 图片处理
- 自动裁剪到标准尺寸
- 添加水印
- 图片格式转换

## 9. 参考资源

- Cloudflare R2 文档：https://developers.cloudflare.com/r2/
- AWS S3 Go SDK：https://aws.github.io/aws-sdk-go-v2/docs/
- Cloudflare Dashboard：https://dash.cloudflare.com/


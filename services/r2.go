package services

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"book-manage/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// R2Service R2服务
type R2Service struct {
	client   *s3.Client
	bucket   string
	publicURL string
}

var r2Service *R2Service

// InitR2Service 初始化R2服务
func InitR2Service(cfg *config.CloudflareR2Config) error {
	// 检查配置是否完整
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.BucketName == "" {
		// 配置不完整，不初始化服务（允许可选）
		return nil
	}

	// 创建S3兼容的配置
	// Cloudflare R2使用自定义端点
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           cfg.Endpoint,
			SigningRegion: cfg.Region,
		}, nil
	})

	cfgOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithEndpointResolverWithOptions(r2Resolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		awsconfig.WithRegion(cfg.Region),
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), cfgOptions...)
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

// GetR2Service 获取R2服务实例
func GetR2Service() *R2Service {
	return r2Service
}

// IsEnabled 检查R2服务是否已启用
func (s *R2Service) IsEnabled() bool {
	return s != nil && s.client != nil
}

// UploadImage 上传图片到R2
func (s *R2Service) UploadImage(bookID int, imageData []byte, filename string) (string, error) {
	if !s.IsEnabled() {
		return "", fmt.Errorf("R2 service is not enabled")
	}

	// 生成唯一文件名
	ext := strings.ToLower(filepath.Ext(filename))
	uniqueFilename := fmt.Sprintf("%d_%s%s", bookID, uuid.New().String(), ext)
	key := fmt.Sprintf("book-covers/%s", uniqueFilename)

	// 上传到R2
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(imageData),
		ContentType: aws.String(getContentType(ext)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	// 返回公开URL
	imageURL := fmt.Sprintf("%s/%s", s.publicURL, key)
	return imageURL, nil
}

// DeleteImage 从R2删除图片
func (s *R2Service) DeleteImage(imageURL string) error {
	if !s.IsEnabled() {
		return fmt.Errorf("R2 service is not enabled")
	}

	// 从URL中提取key
	// 格式：https://pub-xxxxx.r2.dev/book-covers/filename
	key := strings.TrimPrefix(imageURL, s.publicURL+"/")
	if key == imageURL {
		// 如果URL格式不对，尝试其他方式提取
		// 可能URL包含完整路径
		parts := strings.Split(imageURL, "/book-covers/")
		if len(parts) > 1 {
			key = "book-covers/" + parts[1]
		} else {
			return fmt.Errorf("invalid image URL format")
		}
	}

	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

// getContentType 根据文件扩展名获取Content-Type
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


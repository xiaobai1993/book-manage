package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Database   DatabaseConfig   `yaml:"database"`
	Server     ServerConfig     `yaml:"server"`
	JWT        JWTConfig        `yaml:"jwt"`
	Email      EmailConfig      `yaml:"email"`
	AdminEmails []string        `yaml:"admin_emails"`
	CloudflareR2 CloudflareR2Config `yaml:"cloudflare_r2"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `yaml:"secret"`
}

// EmailConfig 邮箱配置
type EmailConfig struct {
	SMTPHost     string `yaml:"smtp_host"`
	SMTPPort     string `yaml:"smtp_port"`
	SMTPUser     string `yaml:"smtp_user"`
	SMTPPassword string `yaml:"smtp_password"`
}

// CloudflareR2Config R2配置（使用S3兼容API）
type CloudflareR2Config struct {
	AccountID       string `yaml:"account_id"`        // 从S3端点URL提取，格式：https://{account-id}.r2.cloudflarestorage.com
	AccessKeyID     string `yaml:"access_key_id"`     // S3 Access Key ID
	SecretAccessKey string `yaml:"secret_access_key"` // S3 Secret Access Key
	BucketName      string `yaml:"bucket_name"`      // 存储桶名称，如：my-object-bucket
	PublicURL       string `yaml:"public_url"`        // 公开访问URL，如：https://pub-xxxxx.r2.dev
	Endpoint        string `yaml:"endpoint"`          // S3端点URL，如：https://{account-id}.r2.cloudflarestorage.com
	Region          string `yaml:"region"`             // 区域，默认：auto
}

// LoadConfig 加载配置
// 环境变量 APP_ENV 可以设置为 env、dev、prod，默认为 env
// 生产环境可以通过环境变量覆盖配置值（优先级：环境变量 > 配置文件）
func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "env"
	}

	// 转换为小写
	env = strings.ToLower(env)

	// 配置文件路径
	configFile := fmt.Sprintf("config/%s.yaml", env)

	// 读取配置文件（如果文件不存在，使用默认值）
	var config Config
	if data, err := os.ReadFile(configFile); err == nil {
		// 解析YAML
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %v", err)
		}
	} else {
		// 文件不存在时使用默认值（PostgreSQL）
		config = Config{
			Database: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "postgres",
				Database: "library_management",
			},
			Server: ServerConfig{
				Port: "8080",
			},
			JWT: JWTConfig{
				Secret: "book-manage-secret-key-2025",
			},
		}
	}

	// 环境变量覆盖（生产环境使用）
	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if database := os.Getenv("DB_NAME"); database != "" {
		config.Database.Database = database
	}

	// 服务器配置
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}

	// JWT 配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWT.Secret = secret
	}

	// 邮箱配置
	if smtpHost := os.Getenv("SMTP_HOST"); smtpHost != "" {
		config.Email.SMTPHost = smtpHost
	}
	if smtpPort := os.Getenv("SMTP_PORT"); smtpPort != "" {
		config.Email.SMTPPort = smtpPort
	}
	if smtpUser := os.Getenv("SMTP_USER"); smtpUser != "" {
		config.Email.SMTPUser = smtpUser
	}
	if smtpPassword := os.Getenv("SMTP_PASSWORD"); smtpPassword != "" {
		config.Email.SMTPPassword = smtpPassword
	}

	// 管理员邮箱（从环境变量读取，用逗号分隔）
	if adminEmails := os.Getenv("ADMIN_EMAILS"); adminEmails != "" {
		config.AdminEmails = strings.Split(adminEmails, ",")
		// 去除空格
		for i, email := range config.AdminEmails {
			config.AdminEmails[i] = strings.TrimSpace(email)
		}
	}

	// Cloudflare R2 配置（从环境变量读取）
	if accountID := os.Getenv("R2_ACCOUNT_ID"); accountID != "" {
		config.CloudflareR2.AccountID = accountID
	}
	if accessKeyID := os.Getenv("R2_ACCESS_KEY_ID"); accessKeyID != "" {
		config.CloudflareR2.AccessKeyID = accessKeyID
	}
	if secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY"); secretAccessKey != "" {
		config.CloudflareR2.SecretAccessKey = secretAccessKey
	}
	if bucketName := os.Getenv("R2_BUCKET_NAME"); bucketName != "" {
		config.CloudflareR2.BucketName = bucketName
	}
	if publicURL := os.Getenv("R2_PUBLIC_URL"); publicURL != "" {
		config.CloudflareR2.PublicURL = publicURL
	}
	if endpoint := os.Getenv("R2_ENDPOINT"); endpoint != "" {
		config.CloudflareR2.Endpoint = endpoint
	}
	if region := os.Getenv("R2_REGION"); region != "" {
		config.CloudflareR2.Region = region
	} else if config.CloudflareR2.Region == "" {
		config.CloudflareR2.Region = "auto" // 默认值
	}

	return &config, nil
}

// IsAdminEmail 检查邮箱是否在管理员白名单中
func (c *Config) IsAdminEmail(email string) bool {
	for _, adminEmail := range c.AdminEmails {
		if strings.ToLower(adminEmail) == strings.ToLower(email) {
			return true
		}
	}
	return false
}

package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Database   DatabaseConfig `yaml:"database"`
	Server     ServerConfig   `yaml:"server"`
	JWT        JWTConfig      `yaml:"jwt"`
	Email      EmailConfig    `yaml:"email"`
	AdminEmails []string      `yaml:"admin_emails"`
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

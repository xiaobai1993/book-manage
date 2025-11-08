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
func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "env"
	}

	// 转换为小写
	env = strings.ToLower(env)

	// 配置文件路径
	configFile := fmt.Sprintf("config/%s.yaml", env)

	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configFile, err)
	}

	// 解析YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
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

package config

// Config 应用配置
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Email    EmailConfig
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string
}

// EmailConfig 邮箱配置
type EmailConfig struct {
	// 这里可以添加真实的邮箱配置，当前版本使用内存存储验证码
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "123456",
			Database: "library_management",
		},
		Server: ServerConfig{
			Port: "8080",
		},
		JWT: JWTConfig{
			Secret: "book-manage-secret-key-2025", // 生产环境应该使用环境变量
		},
		Email: EmailConfig{},
	}
}

package database

import (
	"book-manage/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
// 支持 MySQL 和 PostgreSQL
// 通过环境变量 DB_TYPE 指定数据库类型，默认为 postgres
func InitDB(cfg *config.Config) error {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "postgres" // 默认使用 PostgreSQL
	}
	dbType = strings.ToLower(dbType)

	var err error
	var dsn string

	switch dbType {
	case "postgres", "postgresql":
		// PostgreSQL 连接字符串格式
		// 如果提供了完整的 DATABASE_URL，直接使用（Supabase 会提供）
		if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
			dsn = databaseURL
			// 如果是 Supabase，确保使用 SSL 和正确的参数
			if strings.Contains(databaseURL, "supabase.co") {
				// 如果连接字符串中没有 sslmode 参数，添加它
				if !strings.Contains(databaseURL, "sslmode=") {
					if strings.Contains(databaseURL, "?") {
						dsn = databaseURL + "&sslmode=require"
					} else {
						dsn = databaseURL + "?sslmode=require"
					}
				}
			}
		} else {
			// 否则使用配置文件的参数构建
			// 本地开发环境通常不需要SSL，生产环境（如Supabase）需要SSL
			// 如果host包含supabase.co或不是localhost，则使用require，否则使用disable
			sslmode := "disable"
			if strings.Contains(cfg.Database.Host, "supabase.co") || 
			   strings.Contains(cfg.Database.Host, "amazonaws.com") ||
			   (cfg.Database.Host != "localhost" && cfg.Database.Host != "127.0.0.1") {
				sslmode = "require"
			}
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
				cfg.Database.Host,
				cfg.Database.User,
				cfg.Database.Password,
				cfg.Database.Database,
				cfg.Database.Port,
				sslmode,
			)
		}
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	case "mysql":
		fallthrough
	default:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=10s&writeTimeout=10s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Database,
		)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	// 配置连接池和超时设置
	var sqlDB *sql.DB
	sqlDB, err = DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(25)                 // 最大打开连接数
	sqlDB.SetMaxIdleConns(10)                 // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 空闲连接最大存活时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Database connection established successfully (type: %s)", dbType)
	log.Printf("Connection pool: MaxOpen=%d, MaxIdle=%d", sqlDB.Stats().MaxOpenConnections, sqlDB.Stats().MaxIdleClosed)
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

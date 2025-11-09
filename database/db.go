package database

import (
	"book-manage/config"
	"fmt"
	"log"
	"os"
	"strings"

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
		} else {
			// 否则使用配置文件的参数构建
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
				cfg.Database.Host,
				cfg.Database.User,
				cfg.Database.Password,
				cfg.Database.Database,
				cfg.Database.Port,
			)
		}
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	case "mysql":
		fallthrough
	default:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
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

	log.Printf("Database connection established successfully (type: %s)", dbType)
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

package database

import (
	"book-manage/config"
	"book-manage/models"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
// 仅支持 PostgreSQL
func InitDB(cfg *config.Config) error {
	var err error
	var dsn string

	// 仅支持 PostgreSQL
	{
		// PostgreSQL 连接字符串格式
		// 统一使用 URL 格式：postgresql://user:password@host:port/database?sslmode=xxx
		// 如果提供了完整的 DATABASE_URL，直接使用（Supabase 会提供）
		if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
			dsn = databaseURL
			// 如果是 Supabase 连接池（pooler），尝试转换为直接连接以避免 prepared statement 冲突
			// 连接池的端口是 6543，直接连接的端口是 5432
			// 如果使用连接池遇到 prepared statement 问题，可以设置 USE_DIRECT_CONNECTION=true 环境变量
			useDirectConnection := os.Getenv("USE_DIRECT_CONNECTION") == "true"
			if (strings.Contains(databaseURL, "pooler.supabase.com") ||
				(strings.Contains(databaseURL, "supabase.co") && strings.Contains(databaseURL, ":6543"))) && useDirectConnection {
				// 转换为直接连接：将 pooler.supabase.com 替换为 db.xxx.supabase.co，端口从 6543 改为 5432
				dsn = strings.ReplaceAll(dsn, "pooler.supabase.com", "db."+strings.Split(strings.Split(dsn, "@")[1], ".")[0]+".supabase.co")
				dsn = strings.ReplaceAll(dsn, ":6543", ":5432")
				// 移除用户名中的项目引用（pooler 使用 postgres.xxx，直接连接使用 postgres）
				if strings.Contains(dsn, "postgres.") {
					parts := strings.Split(dsn, "@")
					if len(parts) > 0 {
						userPart := strings.Split(parts[0], "://")[1]
						if strings.Contains(userPart, "postgres.") {
							userPart = strings.ReplaceAll(userPart, "postgres.", "postgres")
							dsn = strings.Split(dsn, "://")[0] + "://" + userPart + "@" + strings.Join(parts[1:], "@")
						}
					}
				}
			}
			// 如果是 Supabase，确保使用 SSL 和正确的参数
			if strings.Contains(dsn, "supabase.co") || strings.Contains(dsn, "pooler.supabase.com") {
				// 如果连接字符串中没有 sslmode 参数，添加它
				if !strings.Contains(dsn, "sslmode=") {
					if strings.Contains(dsn, "?") {
						dsn = dsn + "&sslmode=require"
					} else {
						dsn = dsn + "?sslmode=require"
					}
				}
				// 添加 prefer_simple_protocol=1 来避免 prepared statement 冲突
				// 注意：Supabase 的连接池可能不支持此参数，但 PostgreSQL 14+ 支持
				// 如果连接池不支持，可能需要使用直接连接而不是连接池（设置 USE_DIRECT_CONNECTION=true）
				if !strings.Contains(dsn, "prefer_simple_protocol=") {
					if strings.Contains(dsn, "?") {
						dsn = dsn + "&prefer_simple_protocol=1"
					} else {
						dsn = dsn + "?prefer_simple_protocol=1"
					}
				}
			} else {
				// 非 Supabase 连接也添加 prefer_simple_protocol=1
				if !strings.Contains(dsn, "prefer_simple_protocol=") {
					if strings.Contains(dsn, "?") {
						dsn = dsn + "&prefer_simple_protocol=1"
					} else {
						dsn = dsn + "?prefer_simple_protocol=1"
					}
				}
			}
		} else {
			// 否则使用配置文件的参数构建 URL 格式连接字符串
			// 本地开发环境通常不需要SSL，生产环境（如Supabase）需要SSL
			sslmode := "disable"
			if strings.Contains(cfg.Database.Host, "supabase.co") ||
				strings.Contains(cfg.Database.Host, "pooler.supabase.com") ||
				strings.Contains(cfg.Database.Host, "amazonaws.com") ||
				(cfg.Database.Host != "localhost" && cfg.Database.Host != "127.0.0.1") {
				sslmode = "require"
			}
			// 统一使用 URL 格式构建连接字符串
			// 注意：prefer_simple_protocol 参数仅在 PostgreSQL 14+ 支持
			// 对于 PostgreSQL 12 及以下版本，不添加此参数
			// 本地开发环境通常使用 PostgreSQL 12，生产环境（Supabase）使用 PostgreSQL 14+
			isLocalPostgres := cfg.Database.Host == "localhost" || cfg.Database.Host == "127.0.0.1"
			if isLocalPostgres {
				// 本地环境：不添加 prefer_simple_protocol（兼容 PostgreSQL 12）
				dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=Asia/Shanghai",
					cfg.Database.User,
					cfg.Database.Password,
					cfg.Database.Host,
					cfg.Database.Port,
					cfg.Database.Database,
					sslmode,
				)
			} else {
				// 生产环境（Supabase等）：添加 prefer_simple_protocol（PostgreSQL 14+）
				dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=Asia/Shanghai&prefer_simple_protocol=1",
					cfg.Database.User,
					cfg.Database.Password,
					cfg.Database.Host,
					cfg.Database.Port,
					cfg.Database.Database,
					sslmode,
				)
			}
		}
		// 禁用 prepared statement 以避免 "prepared statement already exists" 错误
		// 这在连接池环境中特别重要，因为连接会被重用
		// 根据 https://github.com/jackc/pgx/issues/1847，这是 pgx 驱动的已知问题
		// 解决方案：使用 postgres.Config 来配置 pgx 驱动，禁用 statement cache
		// PreferSimpleProtocol: true 会强制使用简单协议，完全避免 prepared statement
		postgresConfig := postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // 使用简单协议，完全避免 prepared statement
			WithoutReturning:     false,
		})

		DB, err = gorm.Open(postgresConfig, &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Info),
			PrepareStmt: false, // 禁用 GORM 层面的 prepared statement，避免缓存冲突
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
	// 注意：为了彻底避免 prepared statement 冲突，缩短连接生命周期
	// 这样连接会被更频繁地回收，减少连接重用导致的 prepared statement 冲突
	sqlDB.SetMaxOpenConns(25)                 // 最大打开连接数
	sqlDB.SetMaxIdleConns(5)                  // 最大空闲连接数（减少以降低连接重用）
	sqlDB.SetConnMaxLifetime(2 * time.Minute) // 连接最大生命周期（缩短以减少重用）
	sqlDB.SetConnMaxIdleTime(1 * time.Minute) // 空闲连接最大存活时间（缩短）

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// 自动迁移数据库表
	if err := autoMigrate(); err != nil {
		log.Printf("Warning: Auto migration failed: %v", err)
		// 不阻止启动，因为可能已经手动创建了表
	}

	log.Printf("Database connection established successfully (PostgreSQL)")
	log.Printf("Connection pool: MaxOpen=%d, MaxIdle=%d", sqlDB.Stats().MaxOpenConnections, sqlDB.Stats().MaxIdleClosed)
	return nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// 执行自动迁移（使用models包中的模型）
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.BorrowRecord{},
		&models.EmailCodeRecord{},
	); err != nil {
		return fmt.Errorf("auto migration failed: %v", err)
	}

	log.Printf("Database tables auto-migrated successfully")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

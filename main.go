package main

import (
	"book-manage/config"
	"book-manage/database"
	"book-manage/handlers"
	"book-manage/middleware"
	"book-manage/services"
	"book-manage/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置JWT密钥
	utils.SetJWTSecret(cfg.JWT.Secret)

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化管理员服务
	services.InitAdminService(cfg)

	// 初始化邮件服务
	services.InitEmailService(&cfg.Email)

	// 初始化中间件（传入配置）
	middleware.InitMiddleware(cfg)

	// 创建Gin路由
	r := gin.Default()

	// 添加CORS中间件
	r.Use(middleware.CORSMiddleware())

	// 用户管理模块（无需登录）
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/register", handlers.Register)
		userGroup.POST("/login", handlers.Login)
		userGroup.POST("/sendEmailCode", handlers.SendEmailCode)
		userGroup.POST("/forgetPassword", handlers.ForgetPassword)
	}

	// 用户管理模块（需要登录）
	userAuthGroup := r.Group("/api/user")
	userAuthGroup.Use(middleware.AuthMiddleware())
	{
		userAuthGroup.POST("/profile", handlers.Profile)
		userAuthGroup.POST("/changePassword", handlers.ChangePassword)
		userAuthGroup.POST("/borrowRecords", handlers.BorrowRecords)
	}

	// 图书管理模块（需要登录）
	bookGroup := r.Group("/api/book")
	bookGroup.Use(middleware.AuthMiddleware())
	{
		bookGroup.POST("/detail", handlers.BookDetail)
		bookGroup.POST("/search", handlers.BookSearch)
	}

	// 图书管理模块（需要管理员权限）
	bookAdminGroup := r.Group("/api/book")
	bookAdminGroup.Use(middleware.AuthMiddleware())
	bookAdminGroup.Use(middleware.AdminMiddleware())
	{
		bookAdminGroup.POST("/add", handlers.AddBook)
		bookAdminGroup.POST("/edit", handlers.EditBook)
		bookAdminGroup.POST("/delete", handlers.DeleteBook)
	}

	// 借阅管理模块（需要登录）
	borrowGroup := r.Group("/api/borrow")
	borrowGroup.Use(middleware.AuthMiddleware())
	{
		borrowGroup.POST("/borrow", handlers.Borrow)
		borrowGroup.POST("/return", handlers.Return)
		borrowGroup.POST("/records", handlers.BorrowRecords)
	}

	// 借阅管理模块（需要管理员权限）
	borrowAdminGroup := r.Group("/api/borrow")
	borrowAdminGroup.Use(middleware.AuthMiddleware())
	borrowAdminGroup.Use(middleware.AdminMiddleware())
	{
		borrowAdminGroup.POST("/allRecords", handlers.AllRecords)
	}

	// 管理员模块（需要管理员权限）
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(middleware.AdminMiddleware())
	{
		adminGroup.POST("/emailCodeList", handlers.EmailCodeList)
		adminGroup.POST("/emailCodeStats", handlers.EmailCodeStats)
	}

	// 启动服务器
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	fmt.Printf("API文档请参考 API.md\n")
	fmt.Printf("数据库连接: %s@%s:%s/%s\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

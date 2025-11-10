package handlers

import (
	"book-manage/database"
	"book-manage/models"
	"book-manage/services"
	"book-manage/utils"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddBookRequest 添加图书请求
type AddBookRequest struct {
	Token         string `json:"token"` // token可选，中间件会处理
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`
	ISBN          string `json:"isbn" binding:"required"`
	Category      string `json:"category" binding:"required"`
	TotalQuantity int    `json:"total_quantity" binding:"required"`
	Description   string `json:"description"`
}

// EditBookRequest 编辑图书请求
type EditBookRequest struct {
	Token         string `json:"token"` // token可选，中间件会处理
	ID            int    `json:"id" binding:"required"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	ISBN          string `json:"isbn"`
	Category      string `json:"category"`
	TotalQuantity int    `json:"total_quantity"`
	Description   string `json:"description"`
}

// DeleteBookRequest 删除图书请求
type DeleteBookRequest struct {
	Token string `json:"token"` // token可选，中间件会处理
	ID    int    `json:"id" binding:"required"`
}

// BookDetailRequest 获取图书详情请求
type BookDetailRequest struct {
	Token string `json:"token"` // token可选，中间件会处理
	ID    int    `json:"id" binding:"required"`
}

// BookSearchRequest 图书搜索请求
type BookSearchRequest struct {
	Token    string `json:"token"` // token可选，中间件会处理
	Keyword  string `json:"keyword"`
	Category string `json:"category"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// AddBook 添加图书
func AddBook(c *gin.Context) {
	var req AddBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	db := database.GetDB()

	// 检查ISBN是否已存在
	var existingBook models.Book
	if err := db.Where("isbn = ?", req.ISBN).First(&existingBook).Error; err == nil {
		utils.Error(c, 10017, "ISBN已存在")
		return
	}

	// 创建图书
	book := models.Book{
		Title:             req.Title,
		Author:            req.Author,
		ISBN:              req.ISBN,
		Category:          req.Category,
		TotalQuantity:     req.TotalQuantity,
		AvailableQuantity: req.TotalQuantity,
		Description:       req.Description,
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
	}

	if err := db.Create(&book).Error; err != nil {
		utils.Error(c, 10001, "添加图书失败")
		return
	}

	utils.Success(c, map[string]interface{}{})
}

// EditBook 编辑图书
func EditBook(c *gin.Context) {
	var req EditBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	db := database.GetDB()

	// 查找图书
	var book models.Book
	if err := db.First(&book, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, 10010, "图书不存在")
		} else {
			utils.Error(c, 10001, "查询图书失败")
		}
		return
	}

	// 如果修改了ISBN，检查新ISBN是否已存在
	if req.ISBN != "" && req.ISBN != book.ISBN {
		var existingBook models.Book
		if err := db.Where("isbn = ? AND id != ?", req.ISBN, req.ID).First(&existingBook).Error; err == nil {
			utils.Error(c, 10017, "ISBN已存在")
			return
		}
		book.ISBN = req.ISBN
	}

	// 更新字段
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Category != "" {
		book.Category = req.Category
	}
	if req.Description != "" {
		book.Description = req.Description
	}

	// 如果修改了总数量
	if req.TotalQuantity > 0 {
		// 计算已借出数量
		var borrowedCount int64
		db.Model(&models.BorrowRecord{}).Where("book_id = ? AND status = ?", req.ID, "borrowed").Count(&borrowedCount)

		// 检查总数量是否小于已借出数量
		if req.TotalQuantity < int(borrowedCount) {
			utils.Error(c, 10018, "总数量不能小于已借出数量")
			return
		}

		book.TotalQuantity = req.TotalQuantity
		// 更新可借数量 = 总数量 - 已借出数量
		book.AvailableQuantity = req.TotalQuantity - int(borrowedCount)
	}

	book.UpdateTime = time.Now()

	if err := db.Save(&book).Error; err != nil {
		utils.Error(c, 10001, "更新图书失败")
		return
	}

	utils.Success(c, map[string]interface{}{})
}

// DeleteBook 删除图书
func DeleteBook(c *gin.Context) {
	var req DeleteBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	db := database.GetDB()

	// 检查是否存在未归还的借阅记录
	var borrowCount int64
	db.Model(&models.BorrowRecord{}).Where("book_id = ? AND status = ?", req.ID, "borrowed").Count(&borrowCount)
	if borrowCount > 0 {
		utils.Error(c, 10014, "图书不可删除")
		return
	}

	// 删除图书
	if err := db.Delete(&models.Book{}, req.ID).Error; err != nil {
		utils.Error(c, 10001, "删除图书失败")
		return
	}

	utils.Success(c, map[string]interface{}{})
}

// BookDetail 获取图书详情
func BookDetail(c *gin.Context) {
	var req BookDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	db := database.GetDB()

	// 查找图书
	var book models.Book
	if err := db.First(&book, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, 10010, "图书不存在")
		} else {
			utils.Error(c, 10001, "查询图书失败")
		}
		return
	}

		utils.Success(c, map[string]interface{}{
		"book": map[string]interface{}{
			"id":                book.ID,
			"title":             book.Title,
			"author":            book.Author,
			"isbn":              book.ISBN,
			"category":          book.Category,
			"total_quantity":    book.TotalQuantity,
			"available_quantity": book.AvailableQuantity,
			"description":       book.Description,
			"cover_image_url":   book.CoverImageURL,
			"create_time":       book.CreateTime.Format("2006-01-02 15:04:05"),
			"update_time":       book.UpdateTime.Format("2006-01-02 15:04:05"),
		},
	})
}

// BookSearch 图书搜索
func BookSearch(c *gin.Context) {
	var req BookSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// 验证关键词长度
	if req.Keyword != "" && !utils.ValidateKeyword(req.Keyword) {
		utils.Error(c, 10019, "搜索关键词过短")
		return
	}

	db := database.GetDB()

	// 构建查询
	query := db.Model(&models.Book{})

	// 关键词搜索（书名或作者）
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("title LIKE ? OR author LIKE ?", keyword, keyword)
	}

	// 分类筛选
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 如果没有结果
	if total == 0 && req.Keyword != "" {
		utils.Error(c, 10020, "搜索无结果")
		return
	}

	// 分页查询
	var books []models.Book
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("create_time DESC").Find(&books).Error; err != nil {
		utils.Error(c, 10001, "查询图书失败")
		return
	}

	// 格式化结果
	bookList := make([]map[string]interface{}, len(books))
	for i, book := range books {
		bookList[i] = map[string]interface{}{
			"id":                book.ID,
			"title":             book.Title,
			"author":            book.Author,
			"isbn":              book.ISBN,
			"category":          book.Category,
			"total_quantity":    book.TotalQuantity,
			"available_quantity": book.AvailableQuantity,
			"description":       book.Description,
			"cover_image_url":   book.CoverImageURL,
			"create_time":       book.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	utils.Success(c, map[string]interface{}{
		"books": bookList,
		"total": total,
		"page":  req.Page,
		"limit": req.Limit,
	})
}

// UploadCoverRequest 上传封面请求
type UploadCoverRequest struct {
	BookID int `form:"book_id" binding:"required"`
}

// UploadCover 上传图书封面
func UploadCover(c *gin.Context) {
	var req UploadCoverRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

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

	// 检查R2服务是否可用
	r2Service := services.GetR2Service()
	if !r2Service.IsEnabled() {
		utils.Error(c, 10023, "图片存储服务未配置")
		return
	}

	// 如果已有图片，先删除旧图片
	if book.CoverImageURL != "" {
		if err := r2Service.DeleteImage(book.CoverImageURL); err != nil {
			// 记录错误但不阻止上传新图片
			// log.Printf("Failed to delete old image: %v", err)
		}
	}

	// 上传到R2
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
	BookID int `json:"book_id" binding:"required"`
}

// DeleteCover 删除图书封面
func DeleteCover(c *gin.Context) {
	var req DeleteCoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

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

	// 检查R2服务是否可用
	r2Service := services.GetR2Service()
	if r2Service.IsEnabled() {
		// 从R2删除图片
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

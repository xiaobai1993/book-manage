package handlers

import (
	"book-manage/database"
	"book-manage/models"
	"book-manage/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BorrowRequest 借书请求
type BorrowRequest struct {
	Token  string `json:"token"` // token可选，中间件会处理
	BookID int    `json:"book_id" binding:"required"`
}

// ReturnRequest 还书请求
type ReturnRequest struct {
	Token  string `json:"token"` // token可选，中间件会处理
	BookID int    `json:"book_id" binding:"required"`
}

// BorrowRecordsRequest 获取借阅记录请求
type BorrowRecordsRequest struct {
	Token  string `json:"token"` // token可选，中间件会处理
	Status string `json:"status"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

// AllRecordsRequest 获取全量借阅记录请求
type AllRecordsRequest struct {
	Token      string `json:"token"` // token可选，中间件会处理
	UserEmail  string `json:"user_email"`
	BookTitle  string `json:"book_title"`
	Status     string `json:"status"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

// Borrow 借书
func Borrow(c *gin.Context) {
	var req BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint)

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

	// 检查库存
	if book.AvailableQuantity <= 0 {
		utils.Error(c, 10011, "库存不足")
		return
	}

	// 检查用户当前借阅数量
	var borrowCount int64
	db.Model(&models.BorrowRecord{}).Where("user_id = ? AND status = ?", userIDUint, "borrowed").Count(&borrowCount)
	if borrowCount >= 5 {
		utils.Error(c, 10012, "借阅已达上限")
		return
	}

	// 检查是否已经借阅了该图书且未归还
	var existingRecord models.BorrowRecord
	if err := db.Where("user_id = ? AND book_id = ? AND status = ?", userIDUint, req.BookID, "borrowed").First(&existingRecord).Error; err == nil {
		utils.Error(c, 10013, "该图书已存在未归还记录")
		return
	}

	// 开始事务
	tx := db.Begin()

	// 创建借阅记录
	borrowDate := time.Now()
	dueDate := borrowDate.Add(30 * 24 * time.Hour)
	borrowRecord := models.BorrowRecord{
		UserID:     userIDUint,
		BookID:     uint(req.BookID),
		BorrowDate: borrowDate,
		DueDate:    dueDate,
		Status:     "borrowed",
	}

	if err := tx.Create(&borrowRecord).Error; err != nil {
		tx.Rollback()
		utils.Error(c, 10001, "创建借阅记录失败")
		return
	}

	// 更新图书可借数量
	if err := tx.Model(&book).Update("available_quantity", gorm.Expr("available_quantity - ?", 1)).Error; err != nil {
		tx.Rollback()
		utils.Error(c, 10001, "更新图书库存失败")
		return
	}

	// 提交事务
	tx.Commit()

	utils.Success(c, map[string]interface{}{
		"borrow_date": borrowDate.Format("2006-01-02 15:04:05"),
		"due_date":    dueDate.Format("2006-01-02 15:04:05"),
	})
}

// Return 还书
func Return(c *gin.Context) {
	var req ReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint)

	db := database.GetDB()

	// 查找借阅记录
	var borrowRecord models.BorrowRecord
	if err := db.Where("user_id = ? AND book_id = ? AND status = ?", userIDUint, req.BookID, "borrowed").First(&borrowRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, 10015, "不存在此借阅记录")
		} else {
			utils.Error(c, 10001, "查询借阅记录失败")
		}
		return
	}

	// 检查是否已归还
	if borrowRecord.Status == "returned" {
		utils.Error(c, 10016, "图书已归还")
		return
	}

	// 开始事务
	tx := db.Begin()

	// 更新借阅记录
	returnDate := time.Now()
	if err := tx.Model(&borrowRecord).Updates(map[string]interface{}{
		"status":      "returned",
		"return_date": returnDate,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, 10001, "更新借阅记录失败")
		return
	}

	// 更新图书可借数量
	if err := tx.Model(&models.Book{}).Where("id = ?", req.BookID).Update("available_quantity", gorm.Expr("available_quantity + ?", 1)).Error; err != nil {
		tx.Rollback()
		utils.Error(c, 10001, "更新图书库存失败")
		return
	}

	// 提交事务
	tx.Commit()

	utils.Success(c, map[string]interface{}{})
}

// BorrowRecords 获取借阅记录（个人）
func BorrowRecords(c *gin.Context) {
	var req BorrowRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	db := database.GetDB()

	// 构建查询
	query := db.Table("borrow_record").
		Select("borrow_record.id, borrow_record.book_id, book.title as book_title, borrow_record.borrow_date, borrow_record.due_date, borrow_record.return_date, borrow_record.status").
		Joins("LEFT JOIN book ON borrow_record.book_id = book.id").
		Where("borrow_record.user_id = ?", userIDUint)

	// 状态筛选
	if req.Status != "" && req.Status != "all" {
		query = query.Where("borrow_record.status = ?", req.Status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var records []models.BorrowRecordWithDetails
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("borrow_record.borrow_date DESC").Scan(&records).Error; err != nil {
		utils.Error(c, 10001, "查询借阅记录失败")
		return
	}

	utils.Success(c, map[string]interface{}{
		"records": records,
		"total":   total,
		"page":    req.Page,
		"limit":   req.Limit,
	})
}

// AllRecords 获取全量借阅记录（管理员）
func AllRecords(c *gin.Context) {
	var req AllRecordsRequest
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

	db := database.GetDB()

	// 构建查询
	// 注意：PostgreSQL中user是保留字，需要使用双引号
	query := db.Table("borrow_record").
		Select(`borrow_record.id, borrow_record.user_id, "user".email as user_email, borrow_record.book_id, book.title as book_title, borrow_record.borrow_date, borrow_record.due_date, borrow_record.return_date, borrow_record.status`).
		Joins(`LEFT JOIN "user" ON borrow_record.user_id = "user".id`).
		Joins("LEFT JOIN book ON borrow_record.book_id = book.id")

	// 用户邮箱筛选
	if req.UserEmail != "" {
		query = query.Where(`"user".email LIKE ?`, "%"+req.UserEmail+"%")
	}

	// 图书名称筛选
	if req.BookTitle != "" {
		query = query.Where("book.title LIKE ?", "%"+req.BookTitle+"%")
	}

	// 状态筛选
	if req.Status != "" && req.Status != "all" {
		query = query.Where("borrow_record.status = ?", req.Status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var records []models.BorrowRecordWithDetails
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("borrow_record.borrow_date DESC").Scan(&records).Error; err != nil {
		utils.Error(c, 10001, "查询借阅记录失败")
		return
	}

	utils.Success(c, map[string]interface{}{
		"records": records,
		"total":   total,
		"page":    req.Page,
		"limit":   req.Limit,
	})
}

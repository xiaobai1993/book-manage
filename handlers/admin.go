package handlers

import (
	"book-manage/database"
	"book-manage/models"
	"book-manage/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// EmailCodeListRequest 验证码列表请求
type EmailCodeListRequest struct {
	Page     int    `json:"page" binding:"required,min=1"`
	Limit    int    `json:"limit" binding:"required,min=1,max=100"`
	Email    string `json:"email"`    // 可选：按邮箱筛选
	Action   string `json:"action"`   // 可选：按用途筛选 (register, forget)
	IsUsed   *bool  `json:"is_used"`  // 可选：按是否使用筛选
	Keyword  string `json:"keyword"`  // 可选：关键词搜索（邮箱或验证码）
}

// EmailCodeListResponse 验证码列表响应
type EmailCodeListResponse struct {
	Total int64                  `json:"total"`
	List  []models.EmailCodeRecord `json:"list"`
}

// EmailCodeList 管理员查看验证码记录列表
func EmailCodeList(c *gin.Context) {
	var req EmailCodeListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()
	query := db.Model(&models.EmailCodeRecord{})

	// 按邮箱筛选
	if req.Email != "" {
		query = query.Where("email = ?", req.Email)
	}

	// 按用途筛选
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}

	// 按是否使用筛选
	if req.IsUsed != nil {
		query = query.Where("is_used = ?", *req.IsUsed)
	}

	// 关键词搜索（邮箱或验证码）
	if req.Keyword != "" {
		query = query.Where("email LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var records []models.EmailCodeRecord
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&records).Error; err != nil {
		utils.Error(c, 10001, "查询验证码记录失败")
		return
	}

	utils.Success(c, EmailCodeListResponse{
		Total: total,
		List:  records,
	})
}

// EmailCodeStats 验证码统计信息
func EmailCodeStats(c *gin.Context) {
	db := database.GetDB()

	var stats struct {
		TotalCount    int64 `json:"total_count"`
		UsedCount     int64 `json:"used_count"`
		UnusedCount   int64 `json:"unused_count"`
		ExpiredCount  int64 `json:"expired_count"`
		RegisterCount int64 `json:"register_count"`
		ForgetCount   int64 `json:"forget_count"`
	}

	// 总数
	db.Model(&models.EmailCodeRecord{}).Count(&stats.TotalCount)

	// 已使用数量
	db.Model(&models.EmailCodeRecord{}).Where("is_used = ?", true).Count(&stats.UsedCount)

	// 未使用数量
	db.Model(&models.EmailCodeRecord{}).Where("is_used = ?", false).Count(&stats.UnusedCount)

	// 已过期数量（未使用且已过期）
	now := time.Now()
	db.Model(&models.EmailCodeRecord{}).
		Where("is_used = ? AND expires_at < ?", false, now).
		Count(&stats.ExpiredCount)

	// 注册验证码数量
	db.Model(&models.EmailCodeRecord{}).Where("action = ?", "register").Count(&stats.RegisterCount)

	// 忘记密码验证码数量
	db.Model(&models.EmailCodeRecord{}).Where("action = ?", "forget").Count(&stats.ForgetCount)

	utils.Success(c, stats)
}


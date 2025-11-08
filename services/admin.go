package services

import (
	"book-manage/config"
	"book-manage/database"
	"book-manage/models"
	"strings"

	"gorm.io/gorm"
)

// AdminService 管理员服务
type AdminService struct {
	cfg *config.Config
}

var adminService *AdminService

// InitAdminService 初始化管理员服务
func InitAdminService(cfg *config.Config) {
	adminService = &AdminService{
		cfg: cfg,
	}
}

// GetAdminService 获取管理员服务实例
func GetAdminService() *AdminService {
	return adminService
}

// IsAdmin 判断用户是否为管理员
// 优先级：1. 检查邮箱白名单 2. 检查数据库role字段
func (s *AdminService) IsAdmin(email string) (bool, error) {
	// 优先检查邮箱白名单
	if s.cfg.IsAdminEmail(email) {
		return true, nil
	}

	// 如果不在白名单中，检查数据库role字段
	db := database.GetDB()
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	// 检查role字段是否为admin
	return strings.ToLower(user.Role) == "admin", nil
}

// GetUserRole 获取用户角色
// 返回 "admin" 或 "user"
func (s *AdminService) GetUserRole(email string) (string, error) {
	isAdmin, err := s.IsAdmin(email)
	if err != nil {
		return "user", err
	}

	if isAdmin {
		return "admin", nil
	}
	return "user", nil
}

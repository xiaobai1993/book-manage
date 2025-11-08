package handlers

import (
	"book-manage/database"
	"book-manage/models"
	"book-manage/services"
	"book-manage/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Code            string `json:"code" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SendEmailCodeRequest 发送验证码请求
type SendEmailCodeRequest struct {
	Email  string `json:"email" binding:"required"`
	Action string `json:"action" binding:"required"`
}

// ForgetPasswordRequest 密码找回请求
type ForgetPasswordRequest struct {
	Email              string `json:"email" binding:"required"`
	Code               string `json:"code" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	Token              string `json:"token"` // token可选，中间件会处理
	OldPassword        string `json:"old_password" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
}

// ProfileRequest 获取个人信息请求
type ProfileRequest struct {
	Token string `json:"token"` // token可选，中间件会处理
}


// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	// 验证邮箱格式
	if !utils.ValidateEmail(req.Email) {
		utils.Error(c, 10002, "邮箱格式错误")
		return
	}

	// 验证密码长度
	if !utils.ValidatePassword(req.Password) {
		utils.Error(c, 10006, "密码长度不足")
		return
	}

	// 验证密码一致性
	if req.Password != req.ConfirmPassword {
		utils.Error(c, 10005, "密码不一致")
		return
	}

	// 验证验证码
	emailService := services.GetEmailService()
	if !emailService.VerifyCode(req.Email, "register", req.Code) {
		utils.Error(c, 10004, "验证码错误或已过期")
		return
	}

	db := database.GetDB()

	// 检查邮箱是否已注册
	var existingUser models.User
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.Error(c, 10003, "邮箱已被注册")
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, 10001, "密码加密失败")
		return
	}

	// 创建用户（默认角色为user，可通过配置文件修改）
	user := models.User{
		Email:        req.Email,
		Password:     string(hashedPassword),
		Role:         "user",
		RegisterTime: time.Now(),
		Status:       "normal",
	}

	if err := db.Create(&user).Error; err != nil {
		utils.Error(c, 10001, "注册失败")
		return
	}

	utils.Success(c, map[string]interface{}{})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	db := database.GetDB()

	// 查找用户
	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, 10007, "邮箱或密码错误")
		} else {
			utils.Error(c, 10001, "登录失败")
		}
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(c, 10007, "邮箱或密码错误")
		return
	}

	// 检查账户状态
	if user.Status != "normal" {
		utils.Error(c, 10001, "账户已被禁用")
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		utils.Error(c, 10001, "生成token失败")
		return
	}

	utils.Success(c, map[string]interface{}{
		"user_info": map[string]interface{}{
			"id":            user.ID,
			"email":         user.Email,
			"role":          user.Role,
			"register_time": user.RegisterTime.Format("2006-01-02 15:04:05"),
			"status":        user.Status,
		},
		"token": token,
	})
}

// SendEmailCode 发送邮箱验证码
func SendEmailCode(c *gin.Context) {
	var req SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	// 验证邮箱格式
	if !utils.ValidateEmail(req.Email) {
		utils.Error(c, 10002, "邮箱格式错误")
		return
	}

	// 验证action
	if req.Action != "register" && req.Action != "forget" {
		utils.Error(c, 10001, "action参数错误")
		return
	}

	db := database.GetDB()

	// 如果是注册操作，检查邮箱是否已注册
	if req.Action == "register" {
		var existingUser models.User
		if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			utils.Error(c, 10003, "邮箱已被注册")
			return
		}
	}

	// 如果是忘记密码操作，检查邮箱是否已注册
	if req.Action == "forget" {
		var existingUser models.User
		if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.Error(c, 10008, "邮箱未注册")
			} else {
				utils.Error(c, 10001, "发送验证码失败")
			}
			return
		}
	}

	// 发送验证码
	emailService := services.GetEmailService()
	_, err := emailService.SendCode(req.Email, req.Action)
	if err != nil {
		utils.Error(c, 10001, err.Error())
		return
	}

	utils.Success(c, map[string]interface{}{})
}

// ForgetPassword 密码找回
func ForgetPassword(c *gin.Context) {
	var req ForgetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	// 验证邮箱格式
	if !utils.ValidateEmail(req.Email) {
		utils.Error(c, 10002, "邮箱格式错误")
		return
	}

	// 验证密码长度
	if !utils.ValidatePassword(req.NewPassword) {
		utils.Error(c, 10006, "密码长度不足")
		return
	}

	// 验证密码一致性
	if req.NewPassword != req.ConfirmNewPassword {
		utils.Error(c, 10005, "密码不一致")
		return
	}

	// 验证验证码
	emailService := services.GetEmailService()
	if !emailService.VerifyCode(req.Email, "forget", req.Code) {
		utils.Error(c, 10004, "验证码错误或已过期")
		return
	}

	db := database.GetDB()

	// 查找用户
	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, 10008, "邮箱未注册")
		} else {
			utils.Error(c, 10001, "密码重置失败")
		}
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, 10001, "密码加密失败")
		return
	}

	// 更新密码
	if err := db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		utils.Error(c, 10001, "密码重置失败")
		return
	}

	utils.Success(c, map[string]interface{}{})
}

// Profile 获取个人信息
func Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, 10001, "用户未认证")
		return
	}

	db := database.GetDB()

	// 获取用户信息
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		utils.Error(c, 10001, "获取用户信息失败")
		return
	}

	// 获取当前借阅数量
	var borrowCount int64
	db.Model(&models.BorrowRecord{}).Where("user_id = ? AND status = ?", userID, "borrowed").Count(&borrowCount)

	utils.Success(c, map[string]interface{}{
		"user_info": map[string]interface{}{
			"id":            user.ID,
			"email":         user.Email,
			"role":          user.Role,
			"register_time": user.RegisterTime.Format("2006-01-02 15:04:05"),
			"status":        user.Status,
		},
		"current_borrow_count": borrowCount,
	})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 10001, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")

	// 验证密码长度
	if !utils.ValidatePassword(req.NewPassword) {
		utils.Error(c, 10006, "密码长度不足")
		return
	}

	// 验证密码一致性
	if req.NewPassword != req.ConfirmNewPassword {
		utils.Error(c, 10005, "密码不一致")
		return
	}

	db := database.GetDB()

	// 获取用户信息
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		utils.Error(c, 10001, "获取用户信息失败")
		return
	}

	// 验证原密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.Error(c, 10007, "原密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, 10001, "密码加密失败")
		return
	}

	// 更新密码
	if err := db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		utils.Error(c, 10001, "修改密码失败")
		return
	}

	utils.Success(c, map[string]interface{}{})
}

package services

import (
	"book-manage/config"
	"book-manage/database"
	"book-manage/models"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/resend/resend-go/v2"
)

// EmailCode 邮箱验证码结构
type EmailCode struct {
	Code      string
	Email     string
	Action    string
	ExpiresAt time.Time
}

// EmailService 邮箱服务
type EmailService struct {
	codes  map[string]*EmailCode
	mu     sync.RWMutex
	cfg    *config.EmailConfig
	resend *resend.Client
}

var emailService *EmailService
var once sync.Once

// InitEmailService 初始化邮箱服务
func InitEmailService(cfg *config.EmailConfig) {
	once.Do(func() {
		var client *resend.Client
		if cfg != nil && cfg.SMTPPassword != "" {
			// 使用 Resend API Key 初始化客户端
			client = resend.NewClient(cfg.SMTPPassword)
		}

		emailService = &EmailService{
			codes:  make(map[string]*EmailCode),
			cfg:    cfg,
			resend: client,
		}
		// 启动清理goroutine，定期清理过期验证码
		go emailService.cleanupExpiredCodes()
	})
}

// GetEmailService 获取邮箱服务实例（单例）
func GetEmailService() *EmailService {
	return emailService
}

// GenerateCode 生成6位数字验证码
func (s *EmailService) GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// SendCode 发送验证码
func (s *EmailService) SendCode(email, action string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否在1分钟内重复请求
	key := fmt.Sprintf("%s:%s", email, action)
	if existingCode, exists := s.codes[key]; exists {
		if time.Since(existingCode.ExpiresAt.Add(-29*time.Minute)) < time.Minute {
			return "", fmt.Errorf("请求过于频繁，请稍后再试")
		}
	}

	code := s.GenerateCode()
	expiresAt := time.Now().Add(30 * time.Minute)
	
	// 保存到内存缓存（用于快速验证）
	s.codes[key] = &EmailCode{
		Code:      code,
		Email:     email,
		Action:    action,
		ExpiresAt: expiresAt,
	}

	// 保存到数据库（用于管理员查看）
	db := database.GetDB()
	codeRecord := models.EmailCodeRecord{
		Email:     email,
		Code:      code,
		Action:    action,
		ExpiresAt: expiresAt,
		IsUsed:    false,
	}
	if err := db.Create(&codeRecord).Error; err != nil {
		// 数据库保存失败不影响验证码生成，只记录错误
		fmt.Printf("[Email Service] 保存验证码到数据库失败: %v\n", err)
	}

	// 立即打印验证码（用于调试）
	fmt.Printf("[Email Service] ========== 验证码生成 ==========\n")
	fmt.Printf("[Email Service] 邮箱: %s\n", email)
	fmt.Printf("[Email Service] 操作: %s\n", action)
	fmt.Printf("[Email Service] 验证码: %s\n", code)
	fmt.Printf("[Email Service] 有效期: 30分钟\n")
	fmt.Printf("[Email Service] ================================\n")

	// 发送真实邮件
	if s.resend != nil {
		fmt.Printf("[Email Service] 开始使用 Resend 发送验证码到 %s (action: %s)\n", email, action)
		err := s.sendEmailViaResend(email, action, code)
		if err != nil {
			// 如果发送失败，仍然保留验证码，但记录错误
			fmt.Printf("[Email Service] 发送邮件失败: %v，验证码: %s (用户仍可使用此验证码)\n", err, code)
		}
	} else {
		// 如果未配置邮件服务，打印到控制台（开发模式）
		fmt.Printf("[Email Service] 邮件服务未配置，验证码打印到控制台: %s (%s): %s (有效期30分钟)\n", email, action, code)
	}

	return code, nil
}

// sendEmailViaResend 使用 Resend 发送邮件
func (s *EmailService) sendEmailViaResend(toEmail, action, code string) error {
	if s.resend == nil {
		return fmt.Errorf("Resend 客户端未初始化")
	}

	// 根据action确定邮件主题和内容
	var subject, htmlContent string
	switch action {
	case "register":
		subject = "图书管理系统 - 注册验证码"
		htmlContent = fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #333;">欢迎注册图书管理系统</h2>
				<p>您的注册验证码为：</p>
				<div style="background-color: #f5f5f5; padding: 20px; text-align: center; margin: 20px 0;">
					<h1 style="color: #007bff; font-size: 32px; margin: 0;">%s</h1>
				</div>
				<p>验证码有效期为 30 分钟，请勿泄露给他人。</p>
				<p style="color: #999; font-size: 12px; margin-top: 30px;">此邮件由系统自动发送，请勿回复。</p>
			</div>
		`, code)
	case "forget":
		subject = "图书管理系统 - 密码重置验证码"
		htmlContent = fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #333;">密码重置验证码</h2>
				<p>您正在重置密码，验证码为：</p>
				<div style="background-color: #f5f5f5; padding: 20px; text-align: center; margin: 20px 0;">
					<h1 style="color: #007bff; font-size: 32px; margin: 0;">%s</h1>
				</div>
				<p>验证码有效期为 30 分钟，请勿泄露给他人。</p>
				<p>如非本人操作，请忽略此邮件。</p>
				<p style="color: #999; font-size: 12px; margin-top: 30px;">此邮件由系统自动发送，请勿回复。</p>
			</div>
		`, code)
	default:
		subject = "图书管理系统 - 验证码"
		htmlContent = fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #333;">验证码</h2>
				<p>您的验证码为：</p>
				<div style="background-color: #f5f5f5; padding: 20px; text-align: center; margin: 20px 0;">
					<h1 style="color: #007bff; font-size: 32px; margin: 0;">%s</h1>
				</div>
				<p>验证码有效期为 30 分钟，请勿泄露给他人。</p>
				<p style="color: #999; font-size: 12px; margin-top: 30px;">此邮件由系统自动发送，请勿回复。</p>
			</div>
		`, code)
	}

	// 使用 Resend 发送邮件
	// 发件人邮箱使用 Resend 的域名或者配置的邮箱
	fromEmail := "noreply@yourdomain.com"
	if s.cfg != nil && s.cfg.SMTPUser != "" {
		fromEmail = s.cfg.SMTPUser
	}

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{toEmail},
		Subject: subject,
		Html:    htmlContent,
	}

	sendStart := time.Now()
	fmt.Printf("[Email Service] [Resend] 准备发送邮件到 %s\n", toEmail)

	_, err := s.resend.Emails.Send(params)
	if err != nil {
		fmt.Printf("[Email Service] [Resend] 发送失败 (耗时: %v): %v\n", time.Since(sendStart), err)
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	fmt.Printf("[Email Service] [Resend] 发送成功 (耗时: %v)\n", time.Since(sendStart))
	fmt.Printf("[Email Service] 成功发送验证码邮件到 %s (%s)\n", toEmail, action)
	return nil
}

// VerifyCode 验证验证码
func (s *EmailService) VerifyCode(email, action, code string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s:%s", email, action)
	storedCode, exists := s.codes[key]
	
	// 先从内存缓存验证
	if exists {
		// 检查是否过期
		if time.Now().After(storedCode.ExpiresAt) {
			delete(s.codes, key)
			return false
		}

		// 验证码正确后删除内存缓存
		if storedCode.Code == code {
			delete(s.codes, key)
			
			// 更新数据库记录为已使用
			db := database.GetDB()
			now := time.Now()
			db.Model(&models.EmailCodeRecord{}).
				Where("email = ? AND code = ? AND action = ? AND is_used = ?", email, code, action, false).
				Updates(map[string]interface{}{
					"is_used": true,
					"used_at": &now,
				})
			
			return true
		}
	}

	// 如果内存缓存中没有，尝试从数据库验证（兼容性）
	db := database.GetDB()
	var codeRecord models.EmailCodeRecord
	err := db.Where("email = ? AND code = ? AND action = ? AND is_used = ?", email, code, action, false).
		First(&codeRecord).Error
	
	if err != nil {
		return false
	}

	// 检查是否过期
	if time.Now().After(codeRecord.ExpiresAt) {
		return false
	}

	// 标记为已使用
	now := time.Now()
	db.Model(&codeRecord).Updates(map[string]interface{}{
		"is_used": true,
		"used_at": &now,
	})

	return true
}

// cleanupExpiredCodes 清理过期的验证码
func (s *EmailService) cleanupExpiredCodes() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, code := range s.codes {
			if now.After(code.ExpiresAt) {
				delete(s.codes, key)
			}
		}
		s.mu.Unlock()
	}
}

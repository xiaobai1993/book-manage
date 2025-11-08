package services

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
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
	codes map[string]*EmailCode
	mu    sync.RWMutex
}

var emailService *EmailService
var once sync.Once

// GetEmailService 获取邮箱服务实例（单例）
func GetEmailService() *EmailService {
	once.Do(func() {
		emailService = &EmailService{
			codes: make(map[string]*EmailCode),
		}
		// 启动清理goroutine，定期清理过期验证码
		go emailService.cleanupExpiredCodes()
	})
	return emailService
}

// GenerateCode 生成6位数字验证码
func (s *EmailService) GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// SendCode 发送验证码（当前版本使用内存存储，实际生产环境应发送真实邮件）
func (s *EmailService) SendCode(email, action string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否在1分钟内重复请求
	key := fmt.Sprintf("%s:%s", email, action)
	if existingCode, exists := s.codes[key]; exists {
		if time.Since(existingCode.ExpiresAt.Add(-29 * time.Minute)) < time.Minute {
			return "", fmt.Errorf("请求过于频繁，请稍后再试")
		}
	}

	code := s.GenerateCode()
	s.codes[key] = &EmailCode{
		Code:      code,
		Email:     email,
		Action:    action,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	// 实际生产环境这里应该发送真实邮件
	// 当前版本仅打印到控制台
	fmt.Printf("[Email Service] 发送验证码到 %s (%s): %s (有效期30分钟)\n", email, action, code)

	return code, nil
}

// VerifyCode 验证验证码
func (s *EmailService) VerifyCode(email, action, code string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", email, action)
	storedCode, exists := s.codes[key]
	if !exists {
		return false
	}

	// 检查是否过期
	if time.Now().After(storedCode.ExpiresAt) {
		delete(s.codes, key)
		return false
	}

	// 验证码正确后删除
	if storedCode.Code == code {
		delete(s.codes, key)
		return true
	}

	return false
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

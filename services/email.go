package services

import (
	"book-manage/config"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/smtp"
	"strconv"
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
	cfg   *config.EmailConfig
}

var emailService *EmailService
var once sync.Once

// InitEmailService 初始化邮箱服务
func InitEmailService(cfg *config.EmailConfig) {
	once.Do(func() {
		emailService = &EmailService{
			codes: make(map[string]*EmailCode),
			cfg:   cfg,
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
	s.codes[key] = &EmailCode{
		Code:      code,
		Email:     email,
		Action:    action,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	// 立即打印验证码（用于调试）
	fmt.Printf("[Email Service] ========== 验证码生成 ==========\n")
	fmt.Printf("[Email Service] 邮箱: %s\n", email)
	fmt.Printf("[Email Service] 操作: %s\n", action)
	fmt.Printf("[Email Service] 验证码: %s\n", code)
	fmt.Printf("[Email Service] 有效期: 30分钟\n")
	fmt.Printf("[Email Service] ================================\n")

	// 发送真实邮件
	if s.cfg != nil && s.cfg.SMTPHost != "" && s.cfg.SMTPUser != "" {
		fmt.Printf("[Email Service] 开始发送验证码到 %s (action: %s)\n", email, action)
		err := s.sendEmail(email, action, code)
		if err != nil {
			// 如果发送失败，仍然保留验证码，但记录错误
			fmt.Printf("[Email Service] 发送邮件失败: %v，验证码: %s (用户仍可使用此验证码)\n", err, code)
			// 可以选择返回错误或继续（这里选择继续，至少验证码已生成）
			// return "", fmt.Errorf("发送邮件失败: %v", err)
		}
	} else {
		// 如果未配置邮件服务，打印到控制台（开发模式）
		fmt.Printf("[Email Service] 邮件服务未配置，验证码打印到控制台: %s (%s): %s (有效期30分钟)\n", email, action, code)
	}

	return code, nil
}

// sendEmail 发送邮件
func (s *EmailService) sendEmail(toEmail, action, code string) error {
	if s.cfg == nil {
		return fmt.Errorf("邮件配置未初始化")
	}

	// 打印SMTP配置信息（用于调试，隐藏密码）
	fmt.Printf("[Email Service] ========== SMTP配置信息 ==========\n")
	fmt.Printf("[Email Service] SMTP_HOST: %s\n", s.cfg.SMTPHost)
	fmt.Printf("[Email Service] SMTP_PORT: %s\n", s.cfg.SMTPPort)
	fmt.Printf("[Email Service] SMTP_USER: %s\n", s.cfg.SMTPUser)
	if s.cfg.SMTPPassword != "" {
		// 只显示密码的前2位和后2位，中间用*代替
		pwdLen := len(s.cfg.SMTPPassword)
		if pwdLen > 4 {
			fmt.Printf("[Email Service] SMTP_PASSWORD: %s***%s (长度: %d)\n", 
				s.cfg.SMTPPassword[:2], s.cfg.SMTPPassword[pwdLen-2:], pwdLen)
		} else {
			fmt.Printf("[Email Service] SMTP_PASSWORD: *** (长度: %d)\n", pwdLen)
		}
	} else {
		fmt.Printf("[Email Service] SMTP_PASSWORD: (空)\n")
	}
	fmt.Printf("[Email Service] 收件人: %s\n", toEmail)
	fmt.Printf("[Email Service] 验证码: %s\n", code)
	fmt.Printf("[Email Service] ===================================\n")

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

	// 解析端口
	port, err := strconv.Atoi(s.cfg.SMTPPort)
	if err != nil {
		return fmt.Errorf("无效的SMTP端口: %v", err)
	}

	// 构建邮件内容
	from := s.cfg.SMTPUser
	to := []string{toEmail}

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// 构建邮件消息
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlContent

	// 发送邮件
	sendStart := time.Now()
	if port == 465 {
		// 使用SSL连接（465端口）
		fmt.Printf("[Email Service] 使用SSL方式发送邮件 (端口: %d) 到 %s\n", port, toEmail)
		err = s.sendEmailSSL(s.cfg.SMTPHost, port, s.cfg.SMTPUser, s.cfg.SMTPPassword, from, to, message)
	} else {
		// 使用STARTTLS（587端口等）
		fmt.Printf("[Email Service] 使用STARTTLS方式发送邮件 (端口: %d) 到 %s\n", port, toEmail)
		err = s.sendEmailSTARTTLS(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword, from, to, message)
		if err != nil {
			fmt.Printf("[Email Service] STARTTLS发送失败 (耗时: %v): %v\n", time.Since(sendStart), err)
		} else {
			fmt.Printf("[Email Service] STARTTLS发送成功 (耗时: %v)\n", time.Since(sendStart))
		}
	}

	if err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	fmt.Printf("[Email Service] 成功发送验证码邮件到 %s (%s)\n", toEmail, action)
	return nil
}

// sendEmailSSL 使用SSL发送邮件（465端口）
func (s *EmailService) sendEmailSSL(host string, port int, username, password, from string, to []string, message string) error {
	startTime := time.Now()
	address := fmt.Sprintf("%s:%d", host, port)
	
	fmt.Printf("[Email Service] [SSL] 准备连接到 %s (超时: 5秒)\n", address)

	// 创建带超时的Dialer（5秒超时）
	netDialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	
	// 连接到SMTP服务器（使用5秒超时）
	conn, err := tls.DialWithDialer(netDialer, "tcp", address, &tls.Config{
		ServerName: host,
	})
	if err != nil {
		fmt.Printf("[Email Service] [SSL] 连接失败 (耗时: %v): %v\n", time.Since(startTime), err)
		fmt.Printf("[Email Service] [SSL] 连接地址: %s\n", address)
		fmt.Printf("[Email Service] [SSL] 服务器名称: %s\n", host)
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer conn.Close()
	fmt.Printf("[Email Service] [SSL] 连接成功 (耗时: %v)\n", time.Since(startTime))

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Printf("[Email Service] [SSL] 创建SMTP客户端失败: %v\n", err)
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 认证
	authStart := time.Now()
	auth := smtp.PlainAuth("", username, password, host)
	if err := client.Auth(auth); err != nil {
		fmt.Printf("[Email Service] [SSL] SMTP认证失败 (耗时: %v): %v\n", time.Since(authStart), err)
		return fmt.Errorf("SMTP认证失败: %v", err)
	}
	fmt.Printf("[Email Service] [SSL] SMTP认证成功 (耗时: %v)\n", time.Since(authStart))

	// 设置发件人
	if err := client.Mail(from); err != nil {
		fmt.Printf("[Email Service] [SSL] 设置发件人失败: %v\n", err)
		return fmt.Errorf("设置发件人失败: %v", err)
	}
	fmt.Printf("[Email Service] [SSL] 发件人设置成功: %s\n", from)

	// 设置收件人
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			fmt.Printf("[Email Service] [SSL] 设置收件人失败 (%s): %v\n", recipient, err)
			return fmt.Errorf("设置收件人失败: %v", err)
		}
		fmt.Printf("[Email Service] [SSL] 收件人设置成功: %s\n", recipient)
	}

	// 发送邮件内容
	dataStart := time.Now()
	writer, err := client.Data()
	if err != nil {
		fmt.Printf("[Email Service] [SSL] 准备发送数据失败: %v\n", err)
		return fmt.Errorf("准备发送数据失败: %v", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		writer.Close()
		fmt.Printf("[Email Service] [SSL] 写入邮件内容失败: %v\n", err)
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("[Email Service] [SSL] 关闭数据写入失败: %v\n", err)
		return fmt.Errorf("关闭数据写入失败: %v", err)
	}
	fmt.Printf("[Email Service] [SSL] 邮件内容发送成功 (耗时: %v)\n", time.Since(dataStart))

	quitStart := time.Now()
	err = client.Quit()
	if err != nil {
		fmt.Printf("[Email Service] [SSL] 关闭连接失败: %v\n", err)
		// Quit失败不影响邮件发送，只记录警告
	} else {
		fmt.Printf("[Email Service] [SSL] 连接关闭成功 (耗时: %v)\n", time.Since(quitStart))
	}

	fmt.Printf("[Email Service] [SSL] 邮件发送总耗时: %v\n", time.Since(startTime))
	return nil
}

// sendEmailSTARTTLS 使用STARTTLS发送邮件（587端口等）
func (s *EmailService) sendEmailSTARTTLS(host, port, username, password, from string, to []string, message string) error {
	startTime := time.Now()
	address := fmt.Sprintf("%s:%s", host, port)
	
	fmt.Printf("[Email Service] [STARTTLS] 准备连接到 %s (超时: 5秒)\n", address)

	// 创建带超时的Dialer（5秒超时）
	netDialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	// 先建立TCP连接
	conn, err := netDialer.Dial("tcp", address)
	if err != nil {
		fmt.Printf("[Email Service] [STARTTLS] TCP连接失败 (耗时: %v): %v\n", time.Since(startTime), err)
		fmt.Printf("[Email Service] [STARTTLS] 连接地址: %s\n", address)
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer conn.Close()
	fmt.Printf("[Email Service] [STARTTLS] TCP连接成功 (耗时: %v)\n", time.Since(startTime))

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Printf("[Email Service] [STARTTLS] 创建SMTP客户端失败: %v\n", err)
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 启用STARTTLS
	tlsConfig := &tls.Config{
		ServerName: host,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		fmt.Printf("[Email Service] [STARTTLS] 启用STARTTLS失败: %v\n", err)
		return fmt.Errorf("启用STARTTLS失败: %v", err)
	}
	fmt.Printf("[Email Service] [STARTTLS] STARTTLS启用成功\n")

	// 认证
	authStart := time.Now()
	auth := smtp.PlainAuth("", username, password, host)
	if err := client.Auth(auth); err != nil {
		fmt.Printf("[Email Service] [STARTTLS] SMTP认证失败 (耗时: %v): %v\n", time.Since(authStart), err)
		return fmt.Errorf("SMTP认证失败: %v", err)
	}
	fmt.Printf("[Email Service] [STARTTLS] SMTP认证成功 (耗时: %v)\n", time.Since(authStart))

	// 设置发件人
	if err := client.Mail(from); err != nil {
		fmt.Printf("[Email Service] [STARTTLS] 设置发件人失败: %v\n", err)
		return fmt.Errorf("设置发件人失败: %v", err)
	}
	fmt.Printf("[Email Service] [STARTTLS] 发件人设置成功: %s\n", from)

	// 设置收件人
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			fmt.Printf("[Email Service] [STARTTLS] 设置收件人失败 (%s): %v\n", recipient, err)
			return fmt.Errorf("设置收件人失败: %v", err)
		}
		fmt.Printf("[Email Service] [STARTTLS] 收件人设置成功: %s\n", recipient)
	}

	// 发送邮件内容
	dataStart := time.Now()
	writer, err := client.Data()
	if err != nil {
		fmt.Printf("[Email Service] [STARTTLS] 准备发送数据失败: %v\n", err)
		return fmt.Errorf("准备发送数据失败: %v", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		writer.Close()
		fmt.Printf("[Email Service] [STARTTLS] 写入邮件内容失败: %v\n", err)
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("[Email Service] [STARTTLS] 关闭数据写入失败: %v\n", err)
		return fmt.Errorf("关闭数据写入失败: %v", err)
	}
	fmt.Printf("[Email Service] [STARTTLS] 邮件内容发送成功 (耗时: %v)\n", time.Since(dataStart))

	quitStart := time.Now()
	err = client.Quit()
	if err != nil {
		fmt.Printf("[Email Service] [STARTTLS] 关闭连接失败: %v\n", err)
		// Quit失败不影响邮件发送，只记录警告
	} else {
		fmt.Printf("[Email Service] [STARTTLS] 连接关闭成功 (耗时: %v)\n", time.Since(quitStart))
	}

	fmt.Printf("[Email Service] [STARTTLS] 邮件发送总耗时: %v\n", time.Since(startTime))
	return nil
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

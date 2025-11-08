package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword 验证密码长度
func ValidatePassword(password string) bool {
	return len(password) >= 8
}

// ValidateKeyword 验证搜索关键词长度
func ValidateKeyword(keyword string) bool {
	keyword = strings.TrimSpace(keyword)
	return keyword == "" || len(keyword) >= 2
}

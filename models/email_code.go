package models

import (
	"time"
)

// EmailCodeRecord 邮箱验证码记录模型
type EmailCodeRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Email      string    `gorm:"not null;size:100;index" json:"email"`
	Code       string    `gorm:"not null;size:10" json:"code"`
	Action     string    `gorm:"not null;size:20;index" json:"action"` // register, forget
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;index" json:"created_at"`
	ExpiresAt  time.Time `gorm:"column:expires_at;not null;index" json:"expires_at"`
	IsUsed     bool      `gorm:"column:is_used;default:false;index" json:"is_used"`
	UsedAt     *time.Time `gorm:"column:used_at" json:"used_at"`
}

// TableName 指定表名
func (EmailCodeRecord) TableName() string {
	return "email_code_record"
}


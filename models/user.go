package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Password     string    `gorm:"not null;size:100" json:"-"`
	Role         string    `gorm:"type:enum('admin','user');default:'user'" json:"role"`
	RegisterTime time.Time `gorm:"column:register_time;default:CURRENT_TIMESTAMP" json:"register_time"`
	Status       string    `gorm:"type:enum('normal','disabled');default:'normal'" json:"status"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

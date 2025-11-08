package models

import (
	"time"
)

// Book 图书模型
type Book struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Title            string    `gorm:"not null;size:200" json:"title"`
	Author           string    `gorm:"not null;size:100" json:"author"`
	ISBN             string    `gorm:"uniqueIndex;not null;size:20" json:"isbn"`
	Category         string    `gorm:"not null;size:50" json:"category"`
	TotalQuantity    int       `gorm:"not null" json:"total_quantity"`
	AvailableQuantity int      `gorm:"not null" json:"available_quantity"`
	Description      string    `gorm:"type:text" json:"description"`
	CreateTime       time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 指定表名
func (Book) TableName() string {
	return "book"
}

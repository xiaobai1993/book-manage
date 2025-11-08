package models

import (
	"time"
)

// BorrowRecord 借阅记录模型
type BorrowRecord struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	User       User       `gorm:"foreignKey:UserID" json:"-"`
	BookID     uint       `gorm:"not null;index" json:"book_id"`
	Book       Book       `gorm:"foreignKey:BookID" json:"-"`
	BorrowDate time.Time  `gorm:"column:borrow_date;default:CURRENT_TIMESTAMP" json:"borrow_date"`
	DueDate    time.Time  `gorm:"column:due_date;not null" json:"due_date"`
	ReturnDate *time.Time `gorm:"column:return_date" json:"return_date"`
	Status     string     `gorm:"type:enum('borrowed','returned');default:'borrowed';index" json:"status"`
}

// TableName 指定表名
func (BorrowRecord) TableName() string {
	return "borrow_record"
}

// BorrowRecordWithDetails 包含详细信息的借阅记录
type BorrowRecordWithDetails struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id,omitempty"`
	UserEmail  string    `json:"user_email,omitempty"`
	BookID     uint      `json:"book_id"`
	BookTitle  string    `json:"book_title"`
	BorrowDate time.Time `json:"borrow_date"`
	DueDate    time.Time `json:"due_date"`
	ReturnDate *time.Time `json:"return_date"`
	Status     string    `json:"status"`
}

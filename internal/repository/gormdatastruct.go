package repository

import (
	"time"
)

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
}

type User struct {
	Model
	Login        string `gorm:"unique;not null"`
	HashPassword string `gorm:"type:bytes"`
	Orders       []Order
}

type Order struct {
	Model
	OrderNumber  string `gorm:"unique;not null"`
	UserID       uint
	Transactions []Transaction
	Status       string    `gorm:"size:16"`
	UploadedAt   time.Time `gorm:"autoCreateTime"`
}

type Transaction struct {
	Model
	OrderID         uint    `gorm:"not null"`
	TransactionType string  `gorm:"size:1"`
	Value           float32 `gorm:"type:numeric(8,2)"`
	ProcessedAt     time.Time
}

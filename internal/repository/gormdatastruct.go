package repository

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
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
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type Transaction struct {
	Model
	OrderID         uint    `gorm:"not null"`
	TransactionType string  `gorm:"size:1"`
	Value           float32 `gorm:"type:numeric(8,2)"`
}

type OrderAccrual struct {
	OrderNumber string
	Status      string
	UpdatedAt   time.Time
	Value       float32
	UserID      uint
}

type OrderTransaction struct {
	OrderNumber     string
	Sum             float32
	ProcessedAt     time.Time
	TransactionType string
}

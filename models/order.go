package models

import (
	"time"

	"gorm.io/gorm"
)

type Sell struct {
	ID        uint    `gorm:"primaryKey,autoIncrement" json:"id,omitempty"`
	OrderID   uint    `gorm:"not null" json:"order_id,omitempty"`
	ProductID uint    `gorm:"not null" json:"product_id,omitempty"`
	Product   Product `json:"product,omitempty"`
	Amount    int     `gorm:"default:1" json:"amount,omitempty"`
}

type PurchaseOrder struct {
	ID        uint           `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	RequestID uint           `gorm:"not null" json:"request_id,omitempty"`
	Sells     []Sell         `json:"sells,omitempty"`
}

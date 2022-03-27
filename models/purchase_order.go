package models

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseOrderStatus struct {
	ID              uint      `gorm:"primarykey, autoIncrement" json:"id,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	PurchaseOrderID uint      `json:"purchase_order_id,omitempty" binding:"required"`
	UserID          uint      `json:"user_id,omitempty" binding:"required"`
	User            User      `json:"user,omitempty"`
	Description     string    `json:"description,omitempty" binding:"required"`
}

type Sell struct {
	ID              uint    `gorm:"primaryKey, autoIncrement" json:"id,omitempty"`
	PurchaseOrderID uint    `gorm:"not null" json:"order_id,omitempty" binding:"required"`
	ProductID       uint    `gorm:"not null" json:"product_id,omitempty" binding:"required"`
	Product         Product `json:"product,omitempty"`
	Amount          int     `gorm:"default:1" json:"amount,omitempty" binding:"required"`
}

type PurchaseOrder struct {
	ID        uint                  `gorm:"primarykey" json:"id,omitempty" binding:"required"`
	CreatedAt time.Time             `json:"created_at,omitempty"`
	UpdatedAt time.Time             `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt        `gorm:"index" json:"deleted_at,omitempty"`
	OrderID   uint                  `gorm:"not null" json:"request_id,omitempty" binding:"required"`
	Sells     []Sell                `json:"sells,omitempty" binding:"required"`
	Status    []PurchaseOrderStatus `json:"status,omitempty" binding:"required"`
}

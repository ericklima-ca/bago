package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID             uint            `gorm:"primarykey" json:"id,omitempty" binding:"required"`
	CreatedAt      time.Time       `json:"created_at,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at,omitempty"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
	PurchaseOrders []PurchaseOrder `json:"purchase_orders,omitempty" binding:"required"`
	CustomerID     uint            `json:"customer_id,omitempty" binding:"required"`
	Customer       Customer        `json:"customer,omitempty"`
	CenterID       uint            `json:"center_id,omitempty" binding:"required"`
	Center         Center          `json:"center,omitempty"`
	Status         []OrderStatus   `json:"status,omitempty"`
}

type OrderStatus struct {
	ID          uint      `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	OrderID     uint      `json:"order_id,omitempty" binding:"required"`
	Description string    `gorm:"default:'pending'" json:"description,omitempty" binding:"required"`
	UserID      uint      `json:"user_id,omitempty" binding:"required"`
	User        User      `json:"user,omitempty"`
}

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
}

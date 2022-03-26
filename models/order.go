package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID             uint            `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt      time.Time       `json:"created_at,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at,omitempty"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
	PurchaseOrders []PurchaseOrder `json:"orders,omitempty"`
	CustomerID     uint            `json:"customer_id,omitempty"`
	Customer       Customer        `json:"customer,omitempty"`
	CenterID       uint            `json:"center_id,omitempty"`
	Center         Center          `json:"center,omitempty"`
}

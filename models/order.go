package models

import (
	"time"

	"gorm.io/gorm"
)

type Request struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Orders []Order
	Customer Customer
	Center Center
}

type Product struct {
	ID        uint `gorm:"primarykey"`
	Description string
	ImageUrl string
}

type Order struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	SoldProducts []SoldProduct
}

type SoldProduct struct {
	Product Product
	Amount int
}

type Customer struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name string
	Email string
}

type Center struct {
	ID        uint `gorm:"primarykey"`
	Description string
	Email string
}

type RequestStatus struct {
	Request Request
	Status string
	User string
}
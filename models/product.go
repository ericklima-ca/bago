package models

type Product struct {
	ID          uint    `gorm:"primarykey" json:"id,omitempty"`
	Description string  `json:"description,omitempty"`
	ImageUrl    string  `gorm:"type:TEXT" json:"image_url,omitempty"`
	Price       float64 `gorm:"default:0.0;type:DECIMAL(5,2)" json:"price,omitempty"`
}

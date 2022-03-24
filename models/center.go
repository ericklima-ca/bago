package models

type Center struct {
	ID          uint   `gorm:"primarykey" json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Email       string `json:"email,omitempty"`
}

package models

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type TokenBase struct {
	ID             string `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	ExpirationData time.Time
	UserID         uint `gorm:"primarykey"`
	User           User
}

type TokenSignup struct {
	TokenBase
}

type TokenRecovery struct {
	TokenBase
}

func (u *TokenRecovery) BeforeCreate(*gorm.DB) (err error) {
	u.ExpirationData = time.Now().Add(time.Hour)
	u.ID = generateRandomToken()
	return
}
func (u *TokenSignup) BeforeCreate(*gorm.DB) (err error) {
	u.ExpirationData = time.Now().Add(time.Hour)
	u.ID = generateRandomToken()
	return
}

func generateRandomToken() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 64)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

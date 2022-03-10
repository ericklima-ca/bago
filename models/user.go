package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID             uint `json:"id,omitempty" gorm:"primarykey" binding:"required"`
	CreatedAt      time.Time
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Name           string         `json:"name,omitempty" gorm:"not null" binding:"required"`
	Lastname       string         `json:"lastname,omitempty" gorm:"not null" binding:"required"`
	Email          string         `json:"email,omitempty" gorm:"not null" binding:"required"`
	Role           uint           `json:"role,omitempty"`
	Active         bool           `json:"active" gorm:"default=false"`
	HashedPassword string         `json:"-" gorm:"not null"`
}
type UserFormData struct {
	User
	Password string `json:"password" binding:"required"`
}

func (ufd *UserFormData) GetUser() *User {
	ufd.User.HashedPassword = ufd.Password
	return &ufd.User
}

func (u User) TryAuthenticate(password string) (ok bool) {
	if !u.Active {
		return
	}
	var passToBeCompared = []byte(password)
	var passFromDB = []byte(u.HashedPassword)

	if err := bcrypt.CompareHashAndPassword(passFromDB, passToBeCompared); err == nil {
		ok = true
	}
	return
}

func (u *User) BeforeCreate(*gorm.DB) (err error) {
	userPassBytes := []byte(u.HashedPassword)
	passBytes, err := bcrypt.GenerateFromPassword(userPassBytes, bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("unable to hash password")
	}
	u.HashedPassword = string(passBytes)
	return
}

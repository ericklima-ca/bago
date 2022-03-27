package controllers

import "gorm.io/gorm"

var (
	Auth  AuthController
	Order OrderController
)

func CreateControllers(db *gorm.DB) {
	Auth = AuthController{DB: db}
	Order = OrderController{DB: db}
}

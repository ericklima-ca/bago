package database

import (
	"os"

	"github.com/ericklima-ca/bago/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	Err error
)

func ConnectToDatabase() {

	dsn := os.Getenv("DATABASE_URL")
	DB, Err = gorm.Open(postgres.Open(dsn))
	DB.AutoMigrate(&models.User{}, &models.TokenSignup{}, &models.TokenRecovery{})
}

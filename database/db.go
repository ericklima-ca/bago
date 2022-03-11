package database

import (
	"os"

	"github.com/ericklima-ca/bago/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	Err error
)

func ConnectToDatabase() {

	switch env := os.Getenv("BAGO_ENV"); env {
	case "", "dev":
		dsn := os.Getenv("DATABASE_URL")
		DB, Err = gorm.Open(postgres.Open(dsn))
		DB.AutoMigrate(&models.User{}, &models.TokenSignup{}, &models.TokenRecovery{})
	case "test":
		DB, Err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		DB.AutoMigrate(&models.User{}, &models.TokenSignup{}, &models.TokenRecovery{})
	default:
		panic("No env configured")
	}
}

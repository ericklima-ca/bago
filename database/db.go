package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
type DatabaseServer struct {
	Models []interface{}
	DB  *gorm.DB
	Err error
}

func (dbs *DatabaseServer) Connect() (*gorm.DB, error) {

	switch env := os.Getenv("BAGO_ENV"); env {
	case "", "dev":
		dsn := os.Getenv("DATABASE_URL")
		dbs.DB, dbs.Err = gorm.Open(postgres.Open(dsn))
	case "test":
		dbs.DB, dbs.Err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	default:
		panic("No env configured")
	}
	dbs.DB.AutoMigrate(dbs.Models...)
	return dbs.DB, dbs.Err
}

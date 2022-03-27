package main

import (
	"log"

	"github.com/ericklima-ca/bago/controllers"
	"github.com/ericklima-ca/bago/database"
	"github.com/ericklima-ca/bago/models"
	"github.com/ericklima-ca/bago/router"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	dbs := database.DatabaseServer{
		Models: models.GetModels(),
	}
	db, err := dbs.Connect()
	if err != nil {
		log.Fatalf("DB not connected: %v", err.Error())
	}

	controllers.CreateControllers(db)

	routerServer := router.Router{
		AuthController:  &controllers.Auth,
		OrderController: &controllers.Order,
	}
	server := routerServer.LoadRoutes()

	server.Run(":8080")
}

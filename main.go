package main

import (
	"log"

	"github.com/ericklima-ca/bago/database"
	"github.com/ericklima-ca/bago/router"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	database.ConnectToDatabase()
	if database.Err == nil {
		log.Println("DB connected")
	}
}

func main() {
	routeServer := router.LoadRoutes()

	routeServer.Run(":8080")
}

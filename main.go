package main

import (
	"github.com/ericklima-ca/bago/db"
	"github.com/ericklima-ca/bago/router"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	closer, err := db.ConnectToDatabase()
	if err != nil {
		panic(err)
	}
	defer closer()
}

func main() {
	routeServer := router.LoadRoutes()

	routeServer.Run(":8080")
}

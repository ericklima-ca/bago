package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDatabase() (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, connError := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DATABASE_URL")))
	return func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(fmt.Errorf("error in disconnect %w", err))
		}
	}, connError
}

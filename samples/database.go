package samples

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongo() *mongo.Database {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI("mongodb://localhost:27017"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database("notifications")
}

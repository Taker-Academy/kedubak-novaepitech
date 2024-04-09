// common/db.go

package common

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Client, context.Context) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("'MONGODB_URI' environmental variable not found")
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	return client, context.TODO()
}

func DisconnectFromMongoDB(client *mongo.Client, ctx context.Context) {
	err := client.Disconnect(ctx)

	if err != nil {
		panic(err)
	}
}

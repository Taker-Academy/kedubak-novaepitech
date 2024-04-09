// common/db.go

package common

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertUser(client *mongo.Client, ctx context.Context, user *models.User) {
	collection := client.Database("kedubak").Collection("users")
	user.CreatedAt = time.Now()
	user.LastUpVote = time.Now().Add(-1 * time.Minute)
	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		panic(err)
	}
	objID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		panic("Cannot convert InsertedID to ObjectID")
	}
	user.ID = objID.Hex()
	fmt.Printf("User added with ID %v\n", insertResult.InsertedID)
}


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

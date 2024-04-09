// common/user.go

package common

import (
    "context"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/Taker-Academy/kedubak-novaepitech/models"
)

func GetUserByEmail(client *mongo.Client, ctx context.Context, email string) (*models.User, error) {
    collection := client.Database("kedubak").Collection("users")
    var user models.User
    err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func ComparePasswords(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

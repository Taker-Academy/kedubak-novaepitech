// models/post.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UserID    string             `bson:"userId" json:"userId"`
	FirstName string             `bson:"firstName" json:"firstName"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Comments  []Comment          `bson:"comments" json:"comments"`
	UpVotes   []string           `bson:"upVotes" json:"upVotes"`
}

type Comment struct {
	ID        string `bson:"id" json:"id"`
	FirstName string `bson:"firstName" json:"firstName"`
	Content   string `bson:"content" json:"content"`
}

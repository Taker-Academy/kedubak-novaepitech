package models

import (
	"time"
)

type User struct {
	ID         string    `bson:"_id,omitempty"`
	CreatedAt  time.Time `bson:"createdAt"`
	Email      string    `bson:"email"`
	FirstName  string    `bson:"firstName"`
	LastName   string    `bson:"lastName"`
	Password   string    `bson:"password"`
	LastUpVote time.Time `bson:"lastUpVote"`
}

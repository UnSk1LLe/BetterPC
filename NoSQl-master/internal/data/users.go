package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

var Gmail string

type User struct {
	Name     string    `bson:"name"`
	Surname  string    `bson:"surname"`
	Dob      time.Time `bson:"dob"`
	Email    string    `bson:"email"`
	Password []byte    `bson:"password"`
}

func GetUser(w http.ResponseWriter) User {
	email := Gmail
	var result User
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err := Collection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		http.Error(w, "Error retrieving data from MongoDB", http.StatusInternalServerError)
		return result
	}
	return result
}

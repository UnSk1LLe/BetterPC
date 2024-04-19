package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

var Gmail string

type User struct {
	Name     string
	Surname  string
	Phone    string
	Email    string
	Password string
}

func GetUser(w http.ResponseWriter) User {
	name := Gmail
	var result User
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err := Collection.FindOne(ctx, bson.M{"name": name}).Decode(&result)
	if err != nil {
		http.Error(w, "Error retrieving data from MongoDB", http.StatusInternalServerError)
		return result
	}
	return result
}

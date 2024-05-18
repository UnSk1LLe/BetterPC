package data

import (
	"MongoDb/pkg/logging"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var currentUser User

type User struct {
	ID                primitive.ObjectID `bson:"_id"`
	UserInfo          UserInfo           `bson:"user_info"`
	VerificationToken string             `bson:"verification_token"`
	Verified          bool               `bson:"verified"`
}

type UserInfo struct {
	Name         string    `bson:"name"`
	Surname      string    `bson:"surname"`
	Dob          time.Time `bson:"dob"`
	Email        string    `bson:"email"`
	PasswordHash []byte    `bson:"password"`
}

func CreateUser(user User) error {
	err := Init("test", "users")
	if err != nil {
		return err
	}
	defer CloseConnection()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(email string) (User, error) {
	err := Init("test", "users")
	if err != nil {
		return User{}, err
	}
	defer CloseConnection()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var user User
	err = Collection.FindOne(ctx, bson.M{"user_info.email": email}).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func SetUser(user User) error {
	if user.ID == primitive.NilObjectID {
		logger := logging.GetLogger()
		logger.Infof("ERROR: Trying to set empty user!")
		return errors.New("error: trying to set empty user")
	}
	currentUser = user
	return nil
}

func ClearUser() {
	logger := logging.GetLogger()
	currentUser = User{
		ID: primitive.NilObjectID,
		UserInfo: UserInfo{
			Name:         "",
			Surname:      "",
			Dob:          time.Time{},
			Email:        "",
			PasswordHash: nil,
		},
		VerificationToken: "",
		Verified:          false,
	}
	logger.Infof("User data was cleared")
}

func ShowUser() User {
	if currentUser.ID == primitive.NilObjectID {
		logger := logging.GetLogger()
		logger.Infof("No current user found!")
		return User{}
	}
	return currentUser
}

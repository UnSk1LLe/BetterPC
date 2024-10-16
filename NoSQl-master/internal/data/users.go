package data

import (
	"MongoDb/pkg/session"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
	"time"
)

var UsersCollection *mongo.Collection

type User struct {
	ID                primitive.ObjectID `bson:"_id"`
	UserInfo          UserInfo           `bson:"user_info"`
	SessionToken      string             `bson:"session_token"`
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
	if err != nil || user.ID == primitive.NilObjectID {
		return User{}, err
	}
	return user, nil
}

func GetUserBySessionToken(token string) (User, error) {
	if token == "" {
		return User{}, errors.New("empty session token")
	}
	err := Init("test", "users")
	if err != nil {
		return User{}, err
	}
	defer CloseConnection()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var user User
	err = Collection.FindOne(ctx, bson.M{"session_token": token}).Decode(&user)
	if err != nil || user.ID == primitive.NilObjectID {
		return User{}, err
	}
	return user, nil
}

func GetUserByID(ID primitive.ObjectID) (User, error) {
	if ID == primitive.NilObjectID {
		return User{}, errors.New("user ID cannot be nil")
	}
	err := Init("test", "users")
	if err != nil {
		return User{}, err
	}
	defer CloseConnection()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var user User
	err = Collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func UpdateUser(filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	usersCollection := UsersCollection
	if usersCollection == nil {
		err := Init("test", "users")
		if err != nil {
			return nil, err
		}
		usersCollection = Collection
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ShowUser(r *http.Request, hide ...bool) User {
	user, _ := GetUserBySessionToken(session.GetSessionTokenFromCookie(r))
	if len(hide) > 0 && !hide[0] {
		return user
	}
	user.UserInfo.Email = hideEmail(user.UserInfo.Email)
	return user
}

func IsVerifiedCurrentUser(r *http.Request) bool {
	if !ShowUser(r).Verified {
		return false
	}
	return true
}

func hideEmail(email string) string {
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return email
	}

	localPart := email[:atIndex]
	domainPart := email[atIndex:]

	if len(localPart) <= 2 {

		return localPart + domainPart
	}

	hiddenEmail := localPart[:2] + "***" + domainPart

	return hiddenEmail
}

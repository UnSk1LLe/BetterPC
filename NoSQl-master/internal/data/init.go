package data

import (
	"MongoDb/pkg/logging"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client
var Collection *mongo.Collection
var MongoDbUrl = "mongodb://localhost:27017"

func Init(dbName, collectionName string) error {
	var err error
	logger := logging.GetLogger()
	client, err = mongo.NewClient(options.Client().ApplyURI(MongoDbUrl))
	if err != nil {
		logger.Infof("Could not connect to MongoDB URL: %s", MongoDbUrl)
		log.Fatal(err)
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Infof("Could not connect to MongoDB client: %v", client)
		log.Fatal(err)
		return err
	}

	Collection = client.Database(dbName).Collection(collectionName)
	logger.Infof("Connected to a database: <%s>, collection: <%s>", dbName, collectionName)
	return nil
}

func CloseConnection() {
	logger := logging.GetLogger()

	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			logger.Infof("Error disconnecting from MongoDB: %v", err)
			return
		}
		logger.Infof("MongoDB connection with Client was CLOSED!")
	}
}

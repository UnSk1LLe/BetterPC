package data

import (
	"MongoDb/pkg/logging"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client
var Collection *mongo.Collection
var MongoDbUrl = "mongodb://localhost:27017"

func InitClient() error {
	var err error
	logger := logging.GetLogger()
	client, err = mongo.NewClient(options.Client().ApplyURI(MongoDbUrl))
	if err != nil {
		logger.Fatalf("Could not connect to MongoDB URL: %s", MongoDbUrl)
		log.Fatal(err)
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatalf("Could not connect to MongoDB client: %v", client)
		log.Fatal(err)
		return err
	}
	logger.Infof("connected to a MongoDB URL client: <%v>", MongoDbUrl)
	return nil
}

func Init(dbName, collectionName string) error {
	logger := logging.GetLogger()
	if client == nil {
		err := InitClient()
		if err != nil {
			logger.Fatalf("Error trying to connect to MongoDB client: %v", client)
			return err
		}
	}

	Collection = client.Database(dbName).Collection(collectionName)
	logger.Infof("Connected to a database: <%s>, collection: <%s>", dbName, collectionName)
	return nil
} //TODO use InitCollection() instead

func InitCollection(dbName, collectionName string) (*mongo.Collection, error) {
	logger := logging.GetLogger()
	if client == nil {
		err := InitClient()
		if err != nil {
			logger.Fatalf("Error trying to connect to MongoDB client: %v", client)
			return nil, err
		}
	}

	collection := client.Database(dbName).Collection(collectionName)
	logger.Infof("Connected to a database: <%s>, collection: <%s>", dbName, collectionName)
	return collection, nil
}

func InitAll() {
	logger := logging.GetLogger()
	var err error
	shopDb := "shop"
	userDb := "test"
	if client == nil {
		err = InitClient()
		if err != nil {
			logger.Fatalf("Error trying to connect to MongoDB client: %v", client)
			return
		}
	}
	UsersCollection, err = InitCollection(userDb, "users")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: users", userDb)
	}
	OrdersCollection, err = InitCollection(shopDb, "orders")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: orders", shopDb)
	}
	CpuCollection, err = InitCollection(shopDb, "cpu")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: cpu", shopDb)
	}
	MotherboardCollection, err = InitCollection(shopDb, "motherboard")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: motherboard", shopDb)
	}
	RamCollection, err = InitCollection(shopDb, "ram")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: ram", shopDb)
	}
	SsdCollection, err = InitCollection(shopDb, "ssd")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: ssd", shopDb)
	}
	HddCollection, err = InitCollection(shopDb, "hdd")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: hdd", shopDb)
	}
	GpuCollection, err = InitCollection(shopDb, "gpu")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: gpu", shopDb)
	}
	CoolingCollection, err = InitCollection(shopDb, "cooling")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: cooling", shopDb)
	}
	HousingCollection, err = InitCollection(shopDb, "housing")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: housing", shopDb)
	}
	PowerSupplyCollection, err = InitCollection(shopDb, "powersupply")
	if err != nil {
		logger.Fatalf("Error trying to connect to DB: %v, collection: powersupply", shopDb)
	}
}

func DefineCollection(productType string) (*mongo.Collection, error) {
	logger := logging.GetLogger()
	var collection *mongo.Collection

	switch productType {
	case "cpu":
		collection = CpuCollection
	case "motherboard":
		collection = MotherboardCollection
	case "gpu":
		collection = GpuCollection
	case "cooling":
		collection = CoolingCollection
	case "housing":
		collection = HousingCollection
	case "hdd":
		collection = HddCollection
	case "ssd":
		collection = SsdCollection
	case "ram":
		collection = RamCollection
	case "powersupply":
		collection = PowerSupplyCollection
	default:
		logger.Errorf("wrong product type: %s", productType)
		return nil, fmt.Errorf("wrong product type: %s", productType)
	}
	return collection, nil
}

func CloseConnection() {
	/*logger := logging.GetLogger()

	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			logger.Infof("Error disconnecting from MongoDB: %v", err)
			return
		}
		logger.Infof("MongoDB connection with Client was CLOSED!")
	}*/
	return
}

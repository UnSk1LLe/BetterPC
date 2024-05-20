package data

import (
	"MongoDb/pkg/logging"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var Cart []Item

type Order struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Items  []Item             `bson:"items"`
	UserID primitive.ObjectID `bson:"user_id"`
	Status string             `bson:"status"`
}

type Item struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ProductType string             `bson:"product_type"`
	Model       string             `bson:"model,omitempty"`
	Amount      int                `bson:"amount"`
	Price       int                `bson:"price,omitempty"`
}

func CreateOrder(items []Item, userID primitive.ObjectID) error {
	logger := logging.GetLogger()
	if len(items) == 0 {
		return errors.New("order must contain at least one item")
	} else if userID == primitive.NilObjectID {
		return errors.New("userID cannot be nil")
	}

	newOrder := Order{
		ID:     primitive.NewObjectID(),
		Items:  items,
		UserID: userID,
		Status: "Created",
	}
	err := Init("shop", "orders")
	if err != nil {
		return err
	}
	defer CloseConnection()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = Collection.InsertOne(ctx, newOrder)
	if err != nil {
		return err
	}
	logger.Infof("New Order CREATED: %s", newOrder.ID)
	return nil
}

func GetOrderByID(orderID primitive.ObjectID) (Order, error) {
	logger := logging.GetLogger()
	if orderID == primitive.NilObjectID {
		return Order{}, errors.New("orderID cannot be nil")
	}

	err := Init("shop", "orders")
	if err != nil {
		return Order{}, err
	}
	defer CloseConnection()

	var order Order
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = Collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		return Order{}, err
	}
	logger.Infof("Found Order %v", order)
	return order, nil
}

func GetOrdersByUserID(userID primitive.ObjectID) ([]Order, error) {
	logger := logging.GetLogger()
	if userID == primitive.NilObjectID {
		return nil, errors.New("userID cannot be nil")
	}

	err := Init("shop", "orders")
	if err != nil {
		return nil, err
	}
	defer CloseConnection()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := Collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var orders []Order
	for cursor.Next(ctx) {
		var order Order
		err := cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
	}
	logger.Infof("Found %d Orders", len(orders))
	return orders, nil
}

func SetOrderStatus(orderID primitive.ObjectID, status string) error {
	logger := logging.GetLogger()

	err := Init("shop", "orders")
	if err != nil {
		return err
	}
	defer CloseConnection()

	_, err = Collection.UpdateOne(context.TODO(), bson.M{"_id": orderID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}
	logger.Infof("Set Order #%s STATUS: <%s>", orderID, status)
	return nil
}

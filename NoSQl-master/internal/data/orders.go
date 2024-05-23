package data

import (
	"MongoDb/pkg/logging"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var OrdersCollection *mongo.Collection

type Order struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Items  []Item             `bson:"items"`
	UserID primitive.ObjectID `bson:"user_id"`
	Date   time.Time          `bson:"date"`
	Price  int                `bson:"price"`
	Status string             `bson:"status"`
}

type Item struct {
	ItemHeader   ItemHeader `bson:"item_header"`
	Manufacturer string     `bson:"manufacturer"`
	Model        string     `bson:"model,omitempty"`
	Price        int        `bson:"price,omitempty"`
	MaxAmount    int        `bson:"max_amount,omitempty"`
}

type ItemHeader struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ProductType string             `bson:"product_type"`
	Amount      int                `bson:"amount"`
}

func (i Item) ItemFinalPrice() int {
	finalPrice := i.ItemHeader.Amount * i.Price
	return finalPrice
}

func (o *Order) CalculateOrderPrice() {
	totalPrice := 0
	for _, item := range o.Items {
		totalPrice += item.ItemFinalPrice()
	}
	o.Price = totalPrice
}

func updateProductAmount(productType string, itemID primitive.ObjectID, amountChange int) error {
	collection := defineCollection(productType)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": itemID}
	update := bson.M{"$inc": bson.M{"general.amount": amountChange}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func CreateOrder(items []Item, userID primitive.ObjectID) error {
	logger := logging.GetLogger()
	if len(items) == 0 {
		return errors.New("order must contain at least one item")
	} else if userID == primitive.NilObjectID {
		return errors.New("userID cannot be nil")
	}

	//Check availability
	for _, item := range items {
		collection := defineCollection(item.ItemHeader.ProductType)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var product struct {
			General struct {
				Amount int `bson:"amount"`
			} `bson:"general"`
		}
		err := collection.FindOne(ctx, bson.M{"_id": item.ItemHeader.ID}).Decode(&product)
		if err != nil {
			return err
		}
		if product.General.Amount < item.ItemHeader.Amount {
			return fmt.Errorf("not enough amount for product %s", item.ItemHeader.ID.Hex())
		}

		//Reserve the amount
		err = updateProductAmount(item.ItemHeader.ProductType, item.ItemHeader.ID, -item.ItemHeader.Amount)
		if err != nil {
			return err
		}
	}

	newOrder := Order{
		ID:     primitive.NewObjectID(),
		Items:  items,
		UserID: userID,
		Date:   time.Now(),
		Status: "Created",
	}
	newOrder.CalculateOrderPrice()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := OrdersCollection.InsertOne(ctx, newOrder)
	if err != nil {
		//Rollback the reservation
		for _, item := range items {
			_ = updateProductAmount(item.ItemHeader.ProductType, item.ItemHeader.ID, item.ItemHeader.Amount)
		}
		return err
	}
	logger.Infof("New Order CREATED: %s", newOrder.ID)
	fmt.Println(newOrder)
	return nil
}

func CancelOrder(orderID primitive.ObjectID) error {
	logger := logging.GetLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order Order
	err := OrdersCollection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		return err
	}

	if order.Status == "Canceled" {
		return errors.New("order is already canceled")
	} else if order.Status == "Completed" {
		return errors.New("cannot cancel completed order")
	}

	_, err = OrdersCollection.UpdateOne(ctx, bson.M{"_id": orderID}, bson.M{"$set": bson.M{"status": "Canceled"}})
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		err = updateProductAmount(item.ItemHeader.ProductType, item.ItemHeader.ID, item.ItemHeader.Amount)
		if err != nil {
			return err
		}
	}

	logger.Infof("Order CANCELED: %s", orderID)
	fmt.Println("Order Canceled:", orderID)
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

func GetOrdersByUserID(userID primitive.ObjectID, activeOnly bool) ([]Order, error) {
	logger := logging.GetLogger()
	if userID == primitive.NilObjectID {
		return nil, errors.New("userID cannot be nil")
	}

	err := Init("shop", "orders")
	if err != nil {
		return nil, err
	}
	defer CloseConnection()
	var filter bson.M
	if activeOnly {
		filter = bson.M{
			"user_id": userID,
			"status": bson.M{
				"$nin": []string{"Canceled", "Completed"},
			},
		}
	} else {
		filter = bson.M{"user_id": userID}
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := Collection.Find(ctx, filter)
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
		orders = append(orders, order)
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

	_, err = OrdersCollection.UpdateOne(context.TODO(), bson.M{"_id": orderID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}
	logger.Infof("Set Order #%s STATUS: <%s>", orderID, status)
	return nil
}

func defineCollection(productType string) *mongo.Collection {
	logger := logging.GetLogger()
	var collection *mongo.Collection

	switch productType {
	case "cpu":
		collection = CpuCollection
	case "mbd":
		collection = MotherboardCollection
	case "gpu":
		collection = GpuCollection
	case "clg":
		collection = CoolingCollection
	case "hsg":
		collection = HousingCollection
	case "hdd":
		collection = HddCollection
	case "ssd":
		collection = SsdCollection
	case "ram":
		collection = RamCollection
	case "pwr":
		collection = PowerSupplyCollection
	default:
		logger.Errorf("Error: wrong product type!")
		return nil
	}
	return collection
}

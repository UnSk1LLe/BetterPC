package Handlers

import (
	"MongoDb/internal/data"
	"MongoDb/pkg/logging"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	dbName := "shop"
	var product interface{}
	var collectionName string
	productType := r.FormValue("addToCart")[0:3]

	switch productType {
	case "cpu":
		collectionName = "cpu"
		product = &data.Cpu{}
	case "mbd":
		collectionName = "motherboard"
		product = &data.Motherboard{}
	case "gpu":
		collectionName = "gpu"
		product = &data.Gpu{}
	case "clg":
		collectionName = "cooling"
		product = &data.Cooling{}
	case "hsg":
		collectionName = "housing"
		product = &data.Housing{}
	case "hdd":
		collectionName = "hdd"
		product = &data.Hdd{}
	case "ssd":
		collectionName = "ssd"
		product = &data.Ssd{}
	case "ram":
		collectionName = "ram"
		product = &data.Ram{}
	case "pwr":
		collectionName = "power_supply"
		product = &data.PowerSupply{}
	default:
		logger.Infof("Error: wrong product type!")
		http.Error(w, "Invalid product type", http.StatusBadRequest)
		return
	}

	err := data.Init(dbName, collectionName)
	if err != nil {
		http.Error(w, "Database initialization error", http.StatusInternalServerError)
		return
	}
	defer data.CloseConnection()

	ObjID, err := primitive.ObjectIDFromHex(r.FormValue("addToCart")[13:37])
	if err != nil {
		http.Error(w, "Invalid ObjectID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": ObjID}
	err = data.Collection.FindOne(context.TODO(), filter).Decode(product)
	if err != nil {
		logger.Infof("Error decoding product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	item, err := extractItemFromProduct(product)
	if err != nil {
		logger.Infof("Error extracting item from product: %v", err)
		http.Error(w, "Error extracting product details", http.StatusInternalServerError)
		return
	}

	data.Cart = append(data.Cart, item)
	logger.Infof("Product added to cart: %v", item)

	// Optionally, return a success response or update the UI via JavaScript
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	fmt.Println(data.Cart)
}

func extractItemFromProduct(product interface{}) (data.Item, error) {
	v := reflect.ValueOf(product)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	idField := v.FieldByName("ID")
	generalField := v.FieldByName("General")

	if !idField.IsValid() || !generalField.IsValid() {
		return data.Item{}, errors.New("invalid product structure")
	}

	general := generalField.Interface().(data.General)

	return data.Item{
		ID:     idField.Interface().(primitive.ObjectID),
		Model:  general.Model,
		Price:  general.FinalPrice(),
		Amount: 1,
	}, nil
}

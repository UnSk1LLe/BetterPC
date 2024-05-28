package Handlers

import (
	"MongoDb/internal/data"
	"MongoDb/pkg/logging"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productType := r.FormValue("productType")

	collection, product, err := defineStruct(productType)
	if err != nil {
		HandleError(err, logger, w)
		return
	}

	ObjID, err := primitive.ObjectIDFromHex(r.FormValue("productID")[10:34])
	if err != nil {
		logger.Errorf("Invalid productID: %s", r.FormValue("productID"))
		HandleError(err, logger, w)
		return
	}

	result, err := data.GetProductById(collection, ObjID)
	if err != nil {
		logger.Errorf("Error getting product: %s", r.FormValue("productID"))
		HandleError(err, logger, w)
		return
	}

	err = result.Decode(product)
	if err != nil {
		logger.Errorf("Error decoding product: %v", product)
		HandleError(err, logger, w)
		return
	}

	item, err := extractItemHeaderFromProduct(productType, product)
	if err != nil {
		logger.Errorf("Error extracting item from product: %v", err)
		HandleError(err, logger, w)
		return
	}

	var userCart []data.ItemHeader
	userCart, err = getCartFromCookie(r)
	if err != nil {
		logger.Infof("Error getting cart: %v", err)
		HandleError(err, logger, w)
		return
	}

	for _, cartItem := range userCart {
		if cartItem.ID == item.ID {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
			return
		}
	}

	userCart = append(userCart, item)

	err = saveCartToCookie(w, r, userCart)
	if err != nil {
		logger.Infof("Error saving cart to cookie: %v", err)
		http.Error(w, "Error saving cart", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	logger.Infof("Product added to cart: %v", item)
	fmt.Println(userCart)
}

func extractItemHeaderFromProduct(productType string, product interface{}) (data.ItemHeader, error) {
	v := reflect.ValueOf(product)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	idField := v.FieldByName("ID")

	if !idField.IsValid() {
		return data.ItemHeader{}, errors.New("invalid product structure")
	}

	return data.ItemHeader{
		ID:          idField.Interface().(primitive.ObjectID),
		ProductType: productType,
		Amount:      1,
	}, nil
}

func extractItemFromProduct(itemHeader data.ItemHeader, product interface{}) (data.Item, error) {
	v := reflect.ValueOf(product)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	generalField := v.FieldByName("General")

	if !generalField.IsValid() {
		return data.Item{}, errors.New("invalid product structure")
	}

	general := generalField.Interface().(data.General)

	return data.Item{
		ItemHeader:   itemHeader,
		Manufacturer: general.Manufacturer,
		Model:        general.Model,
		Price:        general.ProductFinalPrice(),
		MaxAmount:    general.Amount,
	}, nil
}

func getCartFromCookie(r *http.Request) ([]data.ItemHeader, error) {
	logger := logging.GetLogger()
	userID := data.ShowUser(r).ID.Hex()
	cartCookie, err := r.Cookie("cart" + userID)

	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return []data.ItemHeader{}, nil
		}
		return nil, err
	}

	//URL-decode the cookie value
	decodedValue, err := url.QueryUnescape(cartCookie.Value)
	if err != nil {
		return nil, err
	}

	fmt.Println("Decoded Cart Cookie Value:", decodedValue)

	var cart []data.ItemHeader
	err = json.Unmarshal([]byte(decodedValue), &cart)
	if err != nil {
		logger.Errorf("Error unmarshalling cart cookie: %v", err)
		logger.Errorf("Cart Cookie Value (Raw): %v", decodedValue)

		return []data.ItemHeader{}, nil
	}

	return cart, nil
}

func saveCartToCookie(w http.ResponseWriter, r *http.Request, cart []data.ItemHeader) error {
	logger := logging.GetLogger()
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	userID := data.ShowUser(r).ID.Hex()

	//URL-encode the JSON string
	encodedValue := url.QueryEscape(string(cartJSON))

	http.SetCookie(w, &http.Cookie{
		Name:    "cart" + userID,
		Value:   encodedValue,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/", //Accessible everywhere
	})

	logger.Infof("Cookie with name: %s was saved.", userID)
	return nil
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	userCart, err := getCartFromCookie(r)
	if err != nil {
		HandleError(errors.New("error getting cart"), logger, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(userCart); err != nil {
		HandleError(errors.New("error encoding cart data"), logger, w)
		return
	}
}

func OpenCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	userCartHeaders, err := getCartFromCookie(r)
	if err != nil {
		logger.Errorf("Error getting cart headers: %v", err)
		HandleError(err, logger, w)
		return
	}

	dbName := "shop"
	userCart, err := getItemsFromItemHeaders(userCartHeaders, dbName, w)
	if err != nil {
		logger.Errorf("Error getting cart data: %v", err)
	}

	dataToSend := struct {
		UserCart []data.Item
		User     data.User
	}{
		UserCart: userCart,
		User:     data.ShowUser(r),
	}

	logger.Infof("Data to send to template: %v", dataToSend)

	tmpl := template.Must(template.ParseFiles("html/cart.html"))
	err = tmpl.Execute(w, dataToSend)
	if err != nil {
		logger.Errorf("Template execution error: %v", err)
		HandleError(err, logger, w)
	}
}

func getItemsFromItemHeaders(userCartHeaders []data.ItemHeader, dbName string, w http.ResponseWriter) ([]data.Item, error) {
	logger := logging.GetLogger()
	var userCart []data.Item

	for _, header := range userCartHeaders {
		collection, product, err := defineStruct(header.ProductType)
		if err != nil {
			HandleError(err, logger, w)
		}

		result, err := data.GetProductById(collection, header.ID)
		if err != nil {
			logger.Errorf("Error fetching product details: %v", err)
			http.Error(w, "Error fetching product details", http.StatusInternalServerError)
			return nil, err
		}
		err = result.Decode(product)
		if err != nil {
			logger.Errorf("Error decoding product: %v", err)
			http.Error(w, "Product not found", http.StatusNotFound)
			return nil, err
		}

		item, err := extractItemFromProduct(header, product)
		userCart = append(userCart, item)
	}
	return userCart, nil
}

func defineStruct(productType string) (*mongo.Collection, interface{}, error) {
	logger := logging.GetLogger()
	var collection *mongo.Collection
	var err error
	var product interface{}

	switch productType {
	case "cpu":
		collection = data.CpuCollection
		product = &data.Cpu{}
	case "motherboard":
		collection = data.MotherboardCollection
		product = &data.Motherboard{}
	case "gpu":
		collection = data.GpuCollection
		product = &data.Gpu{}
	case "cooling":
		collection = data.CoolingCollection
		product = &data.Cooling{}
	case "housing":
		collection = data.HousingCollection
		product = &data.Housing{}
	case "hdd":
		collection = data.HddCollection
		product = &data.Hdd{}
	case "ssd":
		collection = data.SsdCollection
		product = &data.Ssd{}
	case "ram":
		collection = data.RamCollection
		product = &data.Ram{}
	case "powersupply":
		collection = data.PowerSupplyCollection
		product = &data.PowerSupply{}
	default:
		logger.Errorf("Error: wrong product type <%s>!", productType)
		err = errors.New("wrong product type: " + productType)
		return nil, nil, err
	}
	return collection, product, nil
}

func UpdateCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	indexStr := r.FormValue("index")
	quantityStr := r.FormValue("quantity")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	userCart, err := getCartFromCookie(r)
	if err != nil {
		logger.Infof("Error getting cart: %v", err)
		http.Error(w, "Error getting cart", http.StatusInternalServerError)
		return
	}

	if index < 0 || index >= len(userCart) {
		http.Error(w, "Index out of bounds", http.StatusBadRequest)
		return
	}

	userCart[index].Amount = quantity

	err = saveCartToCookie(w, r, userCart)
	if err != nil {
		logger.Infof("Error saving updated cart to cookie: %v", err)
		http.Error(w, "Error saving updated cart", http.StatusInternalServerError)
		return
	}

	logger.Infof("Cart updated: %v", userCart)
}

func DeleteFromCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	itemID, err := primitive.ObjectIDFromHex(r.FormValue("deleteProduct")[10:34])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	userCart, err := getCartFromCookie(r)
	if err != nil {
		logger.Errorf("Error getting cart: %v", err)
		http.Error(w, "Error getting cart", http.StatusInternalServerError)
		return
	}

	var updatedCart []data.ItemHeader
	for _, item := range userCart {
		if item.ID != itemID {
			updatedCart = append(updatedCart, item)
		}
	}

	err = saveCartToCookie(w, r, updatedCart)
	if err != nil {
		logger.Errorf("Error saving updated cart to cookie: %v", err)
		http.Error(w, "Error saving updated cart", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/openCart", http.StatusSeeOther)
	logger.Infof("Product deleted from cart: %v", itemID)
}

func ClearUserCart(w http.ResponseWriter, r *http.Request) {
	userID := data.ShowUser(r).ID.Hex()

	cartCookie, err := r.Cookie("cart" + userID)
	cartCookie.MaxAge = -1
	if err != nil {
		return
	}

	http.SetCookie(w, cartCookie)
}

func ClearUserBuild(w http.ResponseWriter, r *http.Request) {
	userID := data.ShowUser(r).ID.Hex()

	cartCookie, err := r.Cookie("build" + userID)
	cartCookie.MaxAge = -1
	if err != nil {
		return
	}

	http.SetCookie(w, cartCookie)
}

func CreateOrderFromCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	itemHeaders, err := getCartFromCookie(r)
	if err != nil {
		logger.Errorf("Error getting cart from cookie: %v", err)
		HandleError(errors.New("error getting cart from cookie"), logger, w)
		return
	}

	userCart, err := getItemsFromItemHeaders(itemHeaders, "shop", w)
	if err != nil {
		logger.Errorf("Error getting items from itemHeaders: %v", err)
		HandleError(errors.New("error getting items from itemHeaders"), logger, w)
		return
	}

	err = data.CreateOrder(userCart, data.ShowUser(r).ID)
	if err != nil {
		logger.Errorf("Error creating order from cart: %v", err)
		HandleError(errors.New("error creating order from cart"), logger, w)
		return
	}

	ClearUserCart(w, r)
	_ = showMessage("/shop", "Order has been created successfully! You can check it in your profile.", w)
	return
}

func ConvertToItemHeaders(productHeaders []data.ProductHeader) []data.ItemHeader {
	var itemHeaders []data.ItemHeader

	for _, productHeader := range productHeaders {
		itemHeader := data.ItemHeader{
			ID:          productHeader.ID,
			ProductType: productHeader.ProductType,
			Amount:      1,
		}
		itemHeaders = append(itemHeaders, itemHeader)
	}

	return itemHeaders
}

func CreateOrderFromBuild(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productHeaders, err := getBuildFromCookie(r)
	if err != nil {
		logger.Errorf("Error getting build from cookie: %v", err)
		HandleError(errors.New("error getting build from cookie"), logger, w)
		return
	}

	itemHeaders := ConvertToItemHeaders(productHeaders)

	userBuild, err := getItemsFromItemHeaders(itemHeaders, "shop", w)
	if err != nil {
		logger.Errorf("Error getting items from itemHeaders: %v", err)
		HandleError(errors.New("error getting items from itemHeaders"), logger, w)
		return
	}

	err = data.CreateOrder(userBuild, data.ShowUser(r).ID)
	if err != nil {
		logger.Errorf("Error creating order from build: %v", err)
		HandleError(errors.New("error creating order from build"), logger, w)
		return
	}

	ClearUserBuild(w, r)
	_ = showMessage("/shop", "Order has been created successfully! You can check it in your profile.", w)
	return
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	itemID, err := primitive.ObjectIDFromHex(r.FormValue("orderID")[10:34])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	err = data.CancelOrder(itemID)
	if err != nil {
		logger.Errorf("Error canceling order: %v", err)
		http.Error(w, "Error canceling order", http.StatusInternalServerError)
		return
	}
	logger.Infof("Order cancelled: %v", itemID)
	_ = showMessage("/showUserProfile", "Order has been cancelled successfully!", w)
}

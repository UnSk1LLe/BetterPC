package Handlers

import (
	"MongoDb/internal/data"
	"MongoDb/pkg/logging"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	dbName := "shop"
	productType := r.FormValue("addToCart")[0:3]

	collectionName, product := defineStruct(productType)
	if collectionName == "" {
		http.Error(w, "Invalid product type", http.StatusBadRequest)
		return
	}

	ObjID, err := primitive.ObjectIDFromHex(r.FormValue("addToCart")[13:37])
	if err != nil {
		http.Error(w, "Invalid ObjectID", http.StatusBadRequest)
		return
	}

	result, err := data.GetProductById(dbName, collectionName, ObjID)
	if err != nil {
		logger.Infof("Error getting product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	err = result.Decode(product)
	if err != nil {
		logger.Infof("Error decoding product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	item, err := extractItemHeaderFromProduct(productType, product)
	if err != nil {
		logger.Infof("Error extracting item from product: %v", err)
		http.Error(w, "Error extracting product details", http.StatusInternalServerError)
		return
	}

	var userCart []data.ItemHeader
	userCart, err = getCartFromCookie(r)
	if err != nil {
		logger.Infof("Error getting cart: %v", err)
		http.Error(w, "Error getting cart", http.StatusInternalServerError)
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
	userID := data.ShowUser(r).ID.Hex()
	cartCookie, err := r.Cookie(userID)

	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return []data.ItemHeader{}, nil
		}
		return nil, err
	}

	// URL-decode the cookie value
	decodedValue, err := url.QueryUnescape(cartCookie.Value)
	if err != nil {
		return nil, err
	}

	// Debugging output
	fmt.Println("Decoded Cart Cookie Value:", decodedValue)

	var cart []data.ItemHeader
	err = json.Unmarshal([]byte(decodedValue), &cart)
	if err != nil {
		// Additional debugging output
		fmt.Println("Error unmarshalling cart cookie:", err)
		fmt.Println("Cart Cookie Value (Raw):", decodedValue)

		// If the cookie is not a valid JSON, reset it
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

	// URL-encode the JSON string
	encodedValue := url.QueryEscape(string(cartJSON))

	http.SetCookie(w, &http.Cookie{
		Name:    userID,
		Value:   encodedValue,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/", // Accessible everywhere
	})

	logger.Infof("Cookie with name: %s was saved.", userID)
	return nil
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	userCart, err := getCartFromCookie(r)
	if err != nil {
		http.Error(w, "Error getting cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userCart); err != nil {
		http.Error(w, "Error encoding cart data", http.StatusInternalServerError)
		return
	}
}

func OpenCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	// Retrieve cart headers from cookie
	userCartHeaders, err := getCartFromCookie(r)
	if err != nil {
		logger.Infof("Error getting cart headers: %v", err)
		http.Error(w, "Error getting cart", http.StatusInternalServerError)
		return
	}

	dbName := "shop"
	userCart, err := getItemsFromItemHeaders(userCartHeaders, dbName, w)
	if err != nil {
		logger.Errorf("Error getting cart data: %v", err)
		http.Error(w, "Error getting cart", http.StatusInternalServerError)
	}

	dataToSend := []interface{}{userCart, data.ShowUser(r).UserInfo.Name}

	logger.Infof("Data to send to template: %v", dataToSend)

	// Parse and execute the template
	tmpl := template.Must(template.ParseFiles("html/cart.html"))
	err = tmpl.Execute(w, dataToSend)
	if err != nil {
		logger.Infof("Template execution error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func getItemsFromItemHeaders(userCartHeaders []data.ItemHeader, dbName string, w http.ResponseWriter) ([]data.Item, error) {
	logger := logging.GetLogger()
	var userCart []data.Item

	// Fetch full product details for each ItemHeader
	for _, header := range userCartHeaders {
		collectionName, product := defineStruct(header.ProductType)
		if collectionName == "" {
			http.Error(w, "Invalid product type", http.StatusBadRequest)
			continue
		}

		// Fetch product details from the database
		result, err := data.GetProductById(dbName, collectionName, header.ID)
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

func defineStruct(productType string) (string, interface{}) {
	logger := logging.GetLogger()
	var collectionName string
	var product interface{}

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
		return "", nil
	}
	return collectionName, product
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

	cartCookie, err := r.Cookie(userID)
	if err != nil {
		return
	}

	clearCookie := http.Cookie{
		Name:   userID,
		Path:   cartCookie.Path,
		MaxAge: -1,
	}

	http.SetCookie(w, &clearCookie)
}

func CreateOrderFromCart(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	itemHeaders, err := getCartFromCookie(r)
	if err != nil {
		logger.Errorf("Error getting cart from cookie: %v", err)
		http.Error(w, "Error getting cart from cookie", http.StatusInternalServerError)
		return
	}

	userCart, err := getItemsFromItemHeaders(itemHeaders, "shop", w)
	if err != nil {
		logger.Errorf("Error getting items from itemHeaders: %v", err)
		http.Error(w, "Error getting items from itemHeaders", http.StatusInternalServerError)
		return
	}

	err = data.CreateOrder(userCart, data.ShowUser(r).ID)
	if err != nil {
		logger.Errorf("Error creating order from cart: %v", err)
		http.Error(w, "Error creating order from cart", http.StatusInternalServerError)
		return
	}

	ClearUserCart(w, r)
	_ = showMessage("/shop", "Order has been created successfully! You can check it in your profile.", w)
	return
}

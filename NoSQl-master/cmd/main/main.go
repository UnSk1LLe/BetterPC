package main

import (
	"MongoDb/Handlers"
	"MongoDb/internal/data"
	"MongoDb/pkg/logging"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create route")
	data.InitAll()
	http.HandleFunc("/register", Handlers.Register)
	http.HandleFunc("/verify", Handlers.VerifyEmailHandler)
	http.HandleFunc("/recoverPassword", Handlers.RecoverPassword)
	http.HandleFunc("/login", Handlers.Login)
	http.HandleFunc("/logout", Handlers.Logout)
	http.HandleFunc("/home", Handlers.AuthMiddleware(Handlers.Home))
	http.HandleFunc("/shop", Handlers.AuthMiddleware(Handlers.Shop))
	http.HandleFunc("/searchAll", Handlers.AuthMiddleware(Handlers.ListCategories))
	http.HandleFunc("/showUserProfile", Handlers.AuthMiddleware(Handlers.ShowProfile))
	http.HandleFunc("/editUserInfoForm", Handlers.AuthMiddleware(Handlers.EditUserInfoForm))
	http.HandleFunc("/sendVerificationToken", Handlers.AuthMiddleware(Handlers.SendVerificationToken))
	http.HandleFunc("/editUserInfo", Handlers.AuthMiddleware(Handlers.EditUserInfo))
	http.HandleFunc("/showProduct", Handlers.AuthMiddleware(Handlers.ListProductInfo))
	http.HandleFunc("/listProducts", Handlers.AuthMiddleware(Handlers.ListProducts))
	http.HandleFunc("/addProductToBuild", Handlers.AuthMiddleware(Handlers.AddToBuild))
	http.HandleFunc("/deleteProductFromBuild", Handlers.AuthMiddleware(Handlers.DeleteFromBuild))
	http.HandleFunc("/createOrderFromBuild", Handlers.AuthMiddleware(Handlers.CreateOrderFromBuild))

	http.HandleFunc("/getCart", Handlers.AuthMiddleware(Handlers.GetCart))
	http.HandleFunc("/addProductToCart", Handlers.AuthMiddleware(Handlers.AddToCart))
	http.HandleFunc("/openCart", Handlers.AuthMiddleware(Handlers.OpenCart))
	http.HandleFunc("/updateCart", Handlers.AuthMiddleware(Handlers.UpdateCart))
	http.HandleFunc("/deleteProductFromCart", Handlers.AuthMiddleware(Handlers.DeleteFromCart))
	http.HandleFunc("/createOrderFromCart", Handlers.AuthMiddleware(Handlers.CreateOrderFromCart))
	http.HandleFunc("/cancelOrder", Handlers.AuthMiddleware(Handlers.CancelOrder))

	http.HandleFunc("/addNewProductChoice", Handlers.AuthMiddleware(Handlers.AddNewProductChoice))
	http.HandleFunc("/addNewProductForm", Handlers.AuthMiddleware(Handlers.AddNewProductForm))
	http.HandleFunc("/addNewProduct", Handlers.AuthMiddleware(Handlers.AddNewProduct))
	http.HandleFunc("/modifyProductForm", Handlers.AuthMiddleware(Handlers.ModifyProductForm))
	http.HandleFunc("/modifyProduct", Handlers.AuthMiddleware(Handlers.ModifyProduct))
	http.HandleFunc("/deleteProduct", Handlers.AuthMiddleware(Handlers.DeleteProduct))

	/*http.HandleFunc("/productDbTools", Handlers.AuthMiddleware(Handlers.ListProductsDetailed)) //new functions*/

	handler := http.StripPrefix("/assets/", http.FileServer(http.Dir("html/assets")))
	http.Handle("/assets/", handler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logrus.Fatal()
	}
}

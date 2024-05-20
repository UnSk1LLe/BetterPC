package main

import (
	"MongoDb/Handlers"
	"MongoDb/pkg/logging"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create route")
	http.HandleFunc("/register", Handlers.Register)
	http.HandleFunc("/verify", Handlers.VerifyEmailHandler)
	http.HandleFunc("/recoverPassword", Handlers.RecoverPassword)
	http.HandleFunc("/login", Handlers.Login)
	http.HandleFunc("/logout", Handlers.Logout)
	http.HandleFunc("/home", Handlers.AuthMiddleware(Handlers.Home))
	http.HandleFunc("/shop", Handlers.AuthMiddleware(Handlers.Shop))
	http.HandleFunc("/showUserProfile", Handlers.AuthMiddleware(Handlers.ShowProfile))
	http.HandleFunc("/editUserInfoForm", Handlers.AuthMiddleware(Handlers.EditUserInfoForm))
	http.HandleFunc("/editUserInfo", Handlers.AuthMiddleware(Handlers.EditUserInfo))
	http.HandleFunc("/showProduct", Handlers.AuthMiddleware(Handlers.ListProductInfo))
	http.HandleFunc("/listProducts", Handlers.AuthMiddleware(Handlers.ListProducts))
	http.HandleFunc("/cpuFilters", Handlers.AuthMiddleware(Handlers.FilterCpu))
	http.HandleFunc("/addProductToCart", Handlers.AuthMiddleware(Handlers.AddToCart))
	http.HandleFunc("/openCart", Handlers.AuthMiddleware(Handlers.OpenCart))
	http.HandleFunc("/comparisonCpuMb", Handlers.AuthMiddleware(Handlers.ComparisonCpuMb))
	http.HandleFunc("/comparisonCpuRam", Handlers.AuthMiddleware(Handlers.ComparisonCpuRam))
	http.HandleFunc("/comparisonCpuCooling", Handlers.AuthMiddleware(Handlers.ComparisonCpuCooling))
	http.HandleFunc("/comparisonMbCpu", Handlers.AuthMiddleware(Handlers.ComparisonMbCpu))
	http.HandleFunc("/comparisonMbRam", Handlers.AuthMiddleware(Handlers.ComparisonMbRam))
	http.HandleFunc("/comparisonMbHousing", Handlers.AuthMiddleware(Handlers.ComparisonMbHousing))
	http.HandleFunc("/comparisonMbHdd", Handlers.AuthMiddleware(Handlers.ComparisonMbHdd))
	http.HandleFunc("/comparisonMbSsd", Handlers.AuthMiddleware(Handlers.ComparisonMbSsd))
	http.HandleFunc("/comparisonRamCpu", Handlers.AuthMiddleware(Handlers.ComparisonRamCpu))
	http.HandleFunc("/comparisonRamMb", Handlers.AuthMiddleware(Handlers.ComparisonRamMb))
	http.HandleFunc("/comparisonSsdMb", Handlers.AuthMiddleware(Handlers.ComparisonSsdMb))
	http.HandleFunc("/comparisonSsdHousing", Handlers.AuthMiddleware(Handlers.ComparisonSsdHousing))
	http.HandleFunc("/comparisonHddMb", Handlers.AuthMiddleware(Handlers.ComparisonHddMb))
	http.HandleFunc("/comparisonHddHousing", Handlers.AuthMiddleware(Handlers.ComparisonHddHousing))
	http.HandleFunc("/comparisonCoolingCpu", Handlers.AuthMiddleware(Handlers.ComparisonCoolingCpu))
	http.HandleFunc("/addCpuForm", Handlers.AuthMiddleware(Handlers.AddCpuForm))
	http.HandleFunc("/addCpuDo", Handlers.AuthMiddleware(Handlers.AddCpu))
	http.HandleFunc("/modifyCpuForm", Handlers.AuthMiddleware(Handlers.ModifyCpuForm))
	http.HandleFunc("/modifyCpu", Handlers.AuthMiddleware(Handlers.ModifyCpu))
	http.HandleFunc("/deleteCpu", Handlers.AuthMiddleware(Handlers.DeleteCpu))

	handler := http.StripPrefix("/assets/", http.FileServer(http.Dir("html/assets")))
	http.Handle("/assets/", handler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logrus.Fatal()
	}
}

package main

import (
	"MongoDb/Handlers"
	"MongoDb/internal/data"
	"MongoDb/pkg/logging"
	"net/http"
)

func main() {
	logger := logging.GetLogger()
	data.Init("test", "users")
	logger.Info("Create route")
	http.HandleFunc("/register", Handlers.Register)
	http.HandleFunc("/login", Handlers.Login)
	http.HandleFunc("/logout", Handlers.Logout)
	http.HandleFunc("/home", Handlers.AuthMiddleware(Handlers.Home))
	http.HandleFunc("/shop", Handlers.AuthMiddleware(Handlers.Shop))
	http.HandleFunc("/fullListSsd", Handlers.AuthMiddleware(Handlers.ListSsd))
	http.HandleFunc("/fullListCpu", Handlers.AuthMiddleware(Handlers.ListCpu))
	http.HandleFunc("/fullListCooling", Handlers.AuthMiddleware(Handlers.ListCooling))
	http.HandleFunc("/fullListHdd", Handlers.AuthMiddleware(Handlers.ListHdd))
	http.HandleFunc("/fullListHousing", Handlers.AuthMiddleware(Handlers.ListHousing))
	http.HandleFunc("/fullListRam", Handlers.AuthMiddleware(Handlers.ListRam))
	http.HandleFunc("/fullListMotherboard", Handlers.AuthMiddleware(Handlers.ListMotherboard))
	http.HandleFunc("/fullListPowerSupply", Handlers.AuthMiddleware(Handlers.ListPowerSupply))
	http.HandleFunc("/fullListGpu", Handlers.AuthMiddleware(Handlers.ListGpu))
	http.HandleFunc("/listCpu", Handlers.AuthMiddleware(Handlers.List))
	http.HandleFunc("/listCooling", Handlers.AuthMiddleware(Handlers.List))
	http.HandleFunc("/listHdd", Handlers.AuthMiddleware(Handlers.List))
	http.HandleFunc("/listHousing", Handlers.AuthMiddleware(Handlers.List))
	http.HandleFunc("/listRam", Handlers.AuthMiddleware(Handlers.List))
	http.HandleFunc("/listMotherboard", Handlers.AuthMiddleware(Handlers.List))
	http.HandleFunc("/listSsd", Handlers.AuthMiddleware(Handlers.List))
	http.HandleFunc("/listPowersupply", Handlers.AuthMiddleware(Handlers.List))
	http.HandleFunc("/listGpu", Handlers.AuthMiddleware(Handlers.List))
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

	http.ListenAndServe(":8080", nil)
}

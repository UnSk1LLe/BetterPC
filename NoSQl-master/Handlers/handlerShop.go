package Handlers

import (
	"MongoDb/internal/data"
	"MongoDb/pkg/logging"
	"MongoDb/pkg/session"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

//TODO 1 closing connections (1/2 DONE)
//TODO 2 make single function for products list (DONE)
//TODO make single delete, modify functions

func Shop(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/shop.html"))
	dataToSend := struct {
		User data.User
	}{
		User: data.ShowUser(r),
	}
	err := tmpl.Execute(w, dataToSend)
	if err != nil {
		HandleError(err, logging.GetLogger(), w)
		return
	}
}

func ComparisonCpuMb(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listMotherboard.html"))

	data.Init("shop", "motherboard")
	var items []data.Motherboard

	if r.Method == http.MethodPost {

		filter := bson.M{"socket": r.FormValue("mb")}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Motherboard

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		logger.Infof("Found multiple items: %v", len(items))

		dataToSend := []interface{}{items, data.ShowUser(r).UserInfo.Name}

		err = tmpl.Execute(w, dataToSend)
		if err != nil {
			logger.Infof("error: %v", err)
		}
		if err != nil {
			logger.Infof("error :", err)
		}
	}
}

func ComparisonCpuRam(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listRam.html"))

	data.Init("shop", "ram")
	var items []data.Ram

	if r.Method == http.MethodPost {

		filter := bson.M{"type": r.FormValue("ram")}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Ram

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		logger.Infof("Found multiple items: %v", len(items))

		dataToSend := []interface{}{items, data.ShowUser(r).UserInfo.Name}

		err = tmpl.Execute(w, dataToSend)
		if err != nil {
			logger.Infof("error: %v", err)
		}
		if err != nil {
			logger.Infof("error :", err)
		}
	}
}

func ComparisonCpuCooling(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listCooling.html"))

	data.Init("shop", "cooling")
	var items []data.Cooling

	if r.Method == http.MethodPost {

		tdp, _ := strconv.ParseInt(r.FormValue("cooling"), 10, 0)
		filter := bson.M{"tdp": bson.M{"$gte": tdp}}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Cooling

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		logger.Infof("Found multiple items: %v", len(items))

		dataToSend := []interface{}{items, data.ShowUser(r).UserInfo.Name}

		err = tmpl.Execute(w, dataToSend)
		if err != nil {
			logger.Infof("error: %v", err)
		}
		if err != nil {
			logger.Infof("error :", err)
		}
	}
}

func ComparisonMbCpu(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listCpu.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "cpu")
	var items []data.Cpu

	if r.Method == http.MethodPost {

		filter := bson.M{"socket": r.FormValue("cpu")}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Cpu

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonMbRam(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listRam.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "ram")
	var items []data.Ram

	if r.Method == http.MethodPost {

		filter := bson.M{"type": r.FormValue("ram")}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Ram

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonMbHousing(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listHousing.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "housing")
	var items []data.Housing

	if r.Method == http.MethodPost {

		filter := bson.M{"mb_form_factor": r.FormValue("housing")}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Housing

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonMbHdd(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listHdd.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "hdd")
	var items []data.Hdd

	if r.Method == http.MethodPost {

		filter := bson.M{}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Hdd

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonMbSsd(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listSsd.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "ssd")
	var items []data.Ssd

	if r.Method == http.MethodPost {

		filter := bson.M{"interface": "SATA3"}
		pciE := r.FormValue("ssd")
		if pciE == "3" {
			filter = bson.M{"interface": "PCI-E 4x 3.0"}
		} else {
			filter = bson.M{"interface": "PCI-E 4x 4.0"}
		}

		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Ssd

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonRamCpu(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listCpu.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "cpu")
	var items []data.Cpu

	if r.Method == http.MethodPost {

		filter := bson.M{"ram.types": r.FormValue("cpu")}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Cpu

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonRamMb(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listMotherboard.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "motherboard")
	var items []data.Motherboard

	if r.Method == http.MethodPost {

		filter := bson.M{"ram.type": r.FormValue("mb")}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Motherboard

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonSsdMb(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listMotherboard.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "motherboard")
	var items []data.Motherboard

	if r.Method == http.MethodPost {

		var keyTag string

		formFactor := r.FormValue("mb")
		if formFactor != "M2" {
			keyTag = "interfaces.M2"
		} else {
			keyTag = "interfaces.SATA3"
		}

		filter := bson.M{keyTag: bson.M{"$gt": 0}}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Motherboard

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonSsdHousing(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listHousing.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "housing")
	var items []data.Housing

	if r.Method == http.MethodPost {

		filter := bson.M{}
		formFactor := r.FormValue("housing")
		if formFactor != "M2" {
			switch formFactor {
			case "2.5":
				filter = bson.M{"drive_bays." + "2_5": bson.M{"$gte": 0}}
			case "3.5":
				filter = bson.M{"drive_bays." + "3_5": bson.M{"$gte": 0}}
			}
		}

		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Housing

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonHddMb(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listMotherboard.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "motherboard")
	var items []data.Motherboard

	if r.Method == http.MethodPost {
		filter := bson.M{"interfaces.SATA3": bson.M{"$gt": 0}}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Motherboard

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonHddHousing(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listHousing.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "housing")
	var items []data.Housing

	if r.Method == http.MethodPost {

		filter := bson.M{}
		formFactor := r.FormValue("housing")
		if formFactor != "M2" {
			switch formFactor {
			case "2.5":
				filter = bson.M{"drive_bays." + "2_5": bson.M{"$gte": 0}}
			case "3.5":
				filter = bson.M{"drive_bays." + "3_5": bson.M{"$gte": 0}}
			}
		}

		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Housing

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonCoolingCpu(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listCpu.html"))
	tmpl.Execute(w, data.ShowUser(r))

	data.Init("shop", "cpu")
	var items []data.Cpu

	if r.Method == http.MethodPost {

		tdp, _ := strconv.ParseInt(r.FormValue("cpu"), 10, 0)
		filter := bson.M{"tdp": bson.M{"$lte": tdp}}
		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Cpu

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func AddCpuForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/addCpu.html"))
	tmpl.Execute(w, data.ShowUser(r))
}

func AddCpu(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	data.Init("shop", "cpu")

	if r.Method == http.MethodPost {
		year, _ := strconv.Atoi(r.FormValue("year"))
		pcores, _ := strconv.Atoi(r.FormValue("pcores"))
		ecores, _ := strconv.Atoi(r.FormValue("ecores"))
		threads, _ := strconv.Atoi(r.FormValue("threads"))
		techPr, _ := strconv.Atoi(r.FormValue("techPr"))
		pcoresBase, _ := strconv.ParseFloat(r.FormValue("pcoresBase"), 64)
		pcoresBoost, _ := strconv.ParseFloat(r.FormValue("pcoresBoost"), 64)
		ecoresBase, _ := strconv.ParseFloat(r.FormValue("ecoresBase"), 64)
		ecoresBoost, _ := strconv.ParseFloat(r.FormValue("ecoresBoost"), 64)
		if ecores == 0 {
			ecoresBase = 0
			ecoresBoost = 0
		}
		channels, _ := strconv.Atoi(r.FormValue("channels"))
		ramMaxFr, _ := strconv.Atoi(r.FormValue("maxFr"))
		ramMaxCap, _ := strconv.Atoi(r.FormValue("maxCap"))
		tdp, _ := strconv.Atoi(r.FormValue("tdp"))
		pcie, _ := strconv.Atoi(r.FormValue("pcie"))
		maxTemp, _ := strconv.Atoi(r.FormValue("maxTemp"))
		price, _ := strconv.Atoi(r.FormValue("price"))
		discount, _ := strconv.Atoi(r.FormValue("discount"))
		amount, _ := strconv.Atoi(r.FormValue("amount"))
		freeMult := false
		if r.FormValue("freeMult") == "yes" {
			freeMult = true
		}

		general := data.General{
			Manufacturer: r.FormValue("man"),
			Model:        r.FormValue("model"),
			Price:        price,
			Discount:     discount,
			Amount:       amount,
		}

		main := data.MainCpu{
			Category:   r.FormValue("category"),
			Generation: r.FormValue("generation"),
			Socket:     r.FormValue("socket"),
			Year:       year,
		}

		cores := data.CoresCpu{
			Pcores:           pcores,
			Ecores:           ecores,
			Threads:          threads,
			TechnicalProcess: techPr,
		}

		clockFrequency := data.ClockFrequencyCpu{
			Pcores:         []float64{pcoresBase, pcoresBoost},
			Ecores:         []float64{ecoresBase, ecoresBoost},
			FreeMultiplier: freeMult,
		}

		ram := data.RamCpu{
			Channels:     channels,
			Type:         r.FormValue("type"),
			MaxFrequency: ramMaxFr,
			MaxCapacity:  ramMaxCap,
		}

		recordCpu := data.Cpu{
			ID:             primitive.NewObjectID(),
			General:        general,
			Main:           main,
			Cores:          cores,
			ClockFrequency: clockFrequency,
			Ram:            ram,
			Tdp:            tdp,
			Graphics:       r.FormValue("graphics"),
			PciE:           pcie,
			MaxTemperature: maxTemp,
		}

		_, err := data.Collection.InsertOne(context.TODO(), recordCpu)
		if err != nil {
			logger.Infof("A bulk write error occurred: %v", err)
		} else {
			logger.Infof("CPU record with ID: %s was CREATED!", recordCpu.ID)
		}
		http.Redirect(w, r, "/shop", http.StatusSeeOther)
	}
}

func ModifyCpuForm(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/modifyCpu.html"))

	data.Init("shop", "cpu")
	var items []data.Cpu

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("modifyProduct")[13:37])
		filter := bson.M{"_id": ObjID}

		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error: %v", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Cpu

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error: %v", err)
			}

			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error: %v", err)
		}

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error: %v", err)
		}
	}
}

func ModifyCpu(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	err := data.Init("shop", "cpu")
	if err != nil {
		http.Redirect(w, r, "/shop", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		year, _ := strconv.Atoi(r.FormValue("year"))
		pcores, _ := strconv.Atoi(r.FormValue("pcores"))
		ecores, _ := strconv.Atoi(r.FormValue("ecores"))
		threads, _ := strconv.Atoi(r.FormValue("threads"))
		techPr, _ := strconv.Atoi(r.FormValue("techPr"))
		pcoresBase, _ := strconv.ParseFloat(r.FormValue("pcoresBase"), 64)
		pcoresBoost, _ := strconv.ParseFloat(r.FormValue("pcoresBoost"), 64)
		ecoresBase, _ := strconv.ParseFloat(r.FormValue("ecoresBase"), 64)
		ecoresBoost, _ := strconv.ParseFloat(r.FormValue("ecoresBoost"), 64)
		if ecores == 0 {
			ecoresBase = 0
			ecoresBoost = 0
		}
		channels, _ := strconv.Atoi(r.FormValue("channels"))
		ramMaxFr, _ := strconv.Atoi(r.FormValue("maxFr"))
		ramMaxCap, _ := strconv.Atoi(r.FormValue("maxCap"))
		tdp, _ := strconv.Atoi(r.FormValue("tdp"))
		pcie, _ := strconv.Atoi(r.FormValue("pcie"))
		maxTemp, _ := strconv.Atoi(r.FormValue("maxTemp"))
		price, _ := strconv.Atoi(r.FormValue("price"))
		discount, _ := strconv.Atoi(r.FormValue("discount"))
		amount, _ := strconv.Atoi(r.FormValue("amount"))
		freeMult := false
		if r.FormValue("freeMult") == "yes" {
			freeMult = true
		}

		general := data.General{
			Manufacturer: r.FormValue("man"),
			Model:        r.FormValue("model"),
			Price:        price,
			Discount:     discount,
			Amount:       amount,
		}

		main := data.MainCpu{
			Category:   r.FormValue("category"),
			Generation: r.FormValue("generation"),
			Socket:     r.FormValue("socket"),
			Year:       year,
		}

		cores := data.CoresCpu{
			Pcores:           pcores,
			Ecores:           ecores,
			Threads:          threads,
			TechnicalProcess: techPr,
		}

		clockFrequency := data.ClockFrequencyCpu{
			Pcores:         []float64{pcoresBase, pcoresBoost},
			Ecores:         []float64{ecoresBase, ecoresBoost},
			FreeMultiplier: freeMult,
		}

		ram := data.RamCpu{
			Channels:     channels,
			Type:         r.FormValue("type"),
			MaxFrequency: ramMaxFr,
			MaxCapacity:  ramMaxCap,
		}

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("modifyCpu")[10:34])
		filter := bson.M{"_id": ObjID}

		recordCpu := data.Cpu{
			ID:             ObjID,
			General:        general,
			Main:           main,
			Cores:          cores,
			ClockFrequency: clockFrequency,
			Ram:            ram,
			Tdp:            tdp,
			Graphics:       r.FormValue("graphics"),
			PciE:           pcie,
			MaxTemperature: maxTemp,
		}

		update := bson.M{"$set": bson.M{
			"general":         recordCpu.General,
			"main":            recordCpu.Main,
			"cores":           recordCpu.Cores,
			"clock_frequency": recordCpu.ClockFrequency,
			"ram":             recordCpu.Ram,
			"tdp":             recordCpu.Tdp,
			"graphics":        recordCpu.Graphics,
			"pci-e":           recordCpu.PciE,
			"max_temperature": recordCpu.MaxTemperature,
		}}

		_, err = data.Collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			logger.Infof("A bulk write error occurred: %v", err)
		} else {
			logger.Infof("CPU record with ID: %s was UPDATED!", ObjID)
		}

		http.Redirect(w, r, "/shop", http.StatusSeeOther)
	}
}

func DeleteCpu(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	data.Init("shop", "cpu")

	if r.Method == http.MethodPost {
		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("deleteCpu")[10:34])
		filter := bson.M{"_id": ObjID}

		_, err = data.Collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			logger.Infof("A bulk write error occurred: %v", err)
		} else {
			logger.Infof("CPU record with ID: %v was DELETED!", ObjID)
		}
	}
	http.Redirect(w, r, "/shop", http.StatusSeeOther)
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productType := r.URL.Query().Get("productType")
	//listCompatibleOnly := r.URL.Query().Get("listCompatibleOnly")

	var productsList []data.Product
	tmpl := template.Must(template.ParseFiles("html/listProducts.html"))

	filter := getProductsFilter(productType, r)

	productsList, err := productsListing(productType, filter, w)
	if err != nil {
		HandleError(err, logger, w)
		return
	}

	user, _ := data.GetUserBySessionToken(session.GetSessionTokenFromCookie(r))
	build, err := getBuild(r)

	dataToSend := struct {
		ProductType  string
		ProductsList []data.Product
		Build        *data.Build
		User         data.User
	}{
		ProductType:  productType,
		ProductsList: productsList,
		Build:        build,
		User:         user,
	}

	err = tmpl.Execute(w, dataToSend)
	if err != nil {
		HandleError(err, logger, w)
		return
	}
}

func getBuild(r *http.Request) (*data.Build, error) {
	logger := logging.GetLogger()

	productHeaders, err := getBuildFromCookie(r)
	if err != nil {
		logger.Errorf("Failed to list build components: %v", err)
		return nil, err
	}
	//TODO для совместимости получаем FullBuild для отображения получаем Build либо используем только FullBuild а во втором случаем его стандартизируем
	build := &data.Build{}

	getProduct := func(productType string, productID primitive.ObjectID) (data.Product, error) {
		products, err := productsListing(productType, bson.M{"_id": productID}, nil)
		if err != nil || len(products) == 0 {
			return data.Product{}, err
		}
		return products[0], nil
	}

	for _, header := range productHeaders {
		product, err := getProduct(header.ProductType, header.ID)
		if err != nil {
			logger.Errorf("Error getting product: %s", header.ID)
			continue
		}

		switch header.ProductType {
		case "cpu":
			build.CPU = product
		case "motherboard":
			build.Motherboard = product
		case "ram":
			build.RAM = product
		case "gpu":
			build.GPU = product
		case "ssd":
			build.SSD = product
		case "hdd":
			build.HDD = product
		case "cooling":
			build.Cooling = product
		case "powersupply":
			build.PowerSupply = product
		case "housing":
			build.Housing = product
		default:
			logger.Errorf("Unknown product type: %s", header.ProductType)
		}
	}
	return build, nil
}

func productsListing(productType string, filter bson.M, w http.ResponseWriter) ([]data.Product, error) {
	logger := logging.GetLogger()

	dbName := "shop"
	collectionName, _, err := defineStruct(productType)
	if err != nil {
		return nil, err
	}
	var productsList []data.Product

	cur, err := data.GetProducts(dbName, collectionName, filter)
	if err != nil {
		logger.Errorf("Error when trying to get products: %v", err)
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err = cur.Close(ctx)
		if err != nil {
			logger.Errorf("Error when trying to close cursor: %v", err)
			return
		}
	}(cur, context.TODO())

	for cur.Next(context.TODO()) {
		switch productType {
		case "cpu":
			var cpuItem data.Cpu
			err = cur.Decode(&cpuItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			productsList = append(productsList, cpuItem.Standardize())
		case "cooling":
			var coolingItem data.Cooling
			err = cur.Decode(&coolingItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			productsList = append(productsList, coolingItem.Standardize())
		case "motherboard":
			var motherboardItem data.Motherboard
			err = cur.Decode(&motherboardItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			productsList = append(productsList, motherboardItem.Standardize())
		case "housing":
			var housingItem data.Housing
			err = cur.Decode(&housingItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			productsList = append(productsList, housingItem.Standardize())
		case "hdd":
			var hddItem data.Hdd
			err = cur.Decode(&hddItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			productsList = append(productsList, hddItem.Standardize())
		case "ssd":
			var ssdItem data.Ssd
			err = cur.Decode(&ssdItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			productsList = append(productsList, ssdItem.Standardize())
		case "powersupply":
			var powerSupplyItem data.PowerSupply
			err = cur.Decode(&powerSupplyItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			productsList = append(productsList, powerSupplyItem.Standardize())
		case "gpu":
			var gpuItem data.Gpu
			err = cur.Decode(&gpuItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			productsList = append(productsList, gpuItem.Standardize())
		case "ram":
			var ramItem data.Ram
			err = cur.Decode(&ramItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			productsList = append(productsList, ramItem.Standardize())
		default:
			logger.Errorf("Error: invalid product type value!")
			return productsList, errors.New("invalid product type")
		}
	}

	if err = cur.Err(); err != nil {
		logger.Error(err.Error())
	}

	logger.Infof("Found multiple items: %v", len(productsList))
	return productsList, nil
}

func decodeProduct(cur *mongo.Cursor, item interface{}) error {
	logger := logging.GetLogger()
	for cur.Next(context.TODO()) {
		err := cur.Decode(item)
		if err != nil {
			logger.Errorf("error: %v", err)
			continue
		}
		return nil
	}
	return nil
}

func ListProductInfo(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productType := r.FormValue("productType")
	dbName := "shop"

	collectionName, item, err := defineStruct(productType)
	if err != nil {
		HandleError(err, logger, w)
		return
	}

	err = data.Init(dbName, collectionName)
	defer data.CloseConnection()

	ObjID, err := primitive.ObjectIDFromHex(r.FormValue("productID")[10:34])
	filter := bson.M{"_id": ObjID}

	cur, err := data.Collection.Find(context.TODO(), filter)
	if err != nil {
		HandleError(err, logger, w)
		return
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err = cur.Close(ctx)
		if err != nil {
			logger.Errorf("cursor error: %v", err)
		}
	}(cur, context.TODO())

	tmpl := template.Must(template.ParseFiles("html/productInformation.html"))

	err = decodeProduct(cur, item)

	if err != nil {
		HandleError(err, logger, w)
		return
	}

	if err = cur.Err(); err != nil {
		logger.Errorf("error: %v", err)
	}

	logger.Infof("Found Item: %v", item)

	dataToSend := struct {
		ProductType string
		Product     interface{}
		User        data.User
	}{
		ProductType: productType,
		Product:     item,
		User:        data.ShowUser(r),
	}

	err = tmpl.Execute(w, dataToSend)
	if err != nil {
		HandleError(err, logger, w)
		return
	}
}

func appendIntFilterParameters(parameterValues []string, result *[]int) {
	logger := logging.GetLogger()
	for _, value := range parameterValues {
		parameter, err := strconv.Atoi(value)
		if err != nil {
			logger.Errorf("Error parsing core value: %v", err)
			continue
		}
		*result = append(*result, parameter)
	}
}

func appendFloatFilterParameters(parameterValues []string, result *[]float64) {
	logger := logging.GetLogger()
	for _, value := range parameterValues {
		parameter, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.Errorf("Error parsing core value: %v", err)
			continue
		}
		*result = append(*result, parameter)
	}
}

func getInterval(from string, to string) (int, int, error) {
	logger := logging.GetLogger()
	var minValue, maxValue int
	var err error
	if from != "" {
		minValue, err = strconv.Atoi(from)
		if err != nil {
			logger.Errorf("Error parsing price from: %v", err)
			return minValue, maxValue, err
		}
	} else {
		minValue = 0
	}

	if to != "" {
		maxValue, err = strconv.Atoi(to)
		if err != nil {
			logger.Errorf("Error parsing price to: %v", err)
			return minValue, maxValue, err
		}
	} else {
		maxValue = 9999999
	}

	if minValue > maxValue {
		temp := minValue
		minValue = maxValue
		maxValue = temp
	}
	return minValue, maxValue, nil
}

func filterCpu(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var categories []string
	var cores []int
	var threads []int
	var priceFrom, priceTo int
	var ramTypes []string
	var sockets []string
	var pcie []int

	manufacturers = r.Form["Manufacturer"]
	categories = r.Form["Category"]
	coreValues := r.Form["Cores"]
	appendIntFilterParameters(coreValues, &cores)

	threadsValues := r.Form["Threads"]
	appendIntFilterParameters(threadsValues, &threads)

	priceFromStr := r.Form.Get("Price-min")
	priceToStr := r.Form.Get("Price-max")

	priceFrom, priceTo, err = getInterval(priceFromStr, priceToStr)

	ramTypes = r.Form["ram-type"]
	sockets = r.Form["socket"]
	pcieValues := r.Form["pcie-version"]
	appendIntFilterParameters(pcieValues, &pcie)

	filter := bson.M{}

	integratedGraphicsValues := r.Form["integrated-graphics"]
	var integrated bson.M

	if len(integratedGraphicsValues) == 1 {
		value := integratedGraphicsValues[0]
		if value == "yes" {
			integrated = bson.M{"$ne": ""}
		} else if value == "no" {
			integrated = bson.M{"$eq": ""}
		}
		filter["graphics"] = integrated
	}

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(categories) > 0 {
		categoryFilter := bson.M{"$in": categories}
		filter["main.category"] = categoryFilter
	}

	if len(cores) > 0 {
		filter = bson.M{
			"$expr": bson.M{
				"$in": bson.A{
					bson.M{"$add": bson.A{"$cores.p-cores", "$cores.e-cores"}},
					cores,
				},
			},
		}
	}

	if len(threads) > 0 {
		filter["cores.threads"] = bson.M{"$in": threads}
	}

	if priceFrom != 0 || priceTo != 0 {
		filter["general.price"] = bson.M{"$gte": priceFrom, "$lte": priceTo}
	}

	if len(ramTypes) > 0 {
		filter["ram.type"] = bson.M{"$in": ramTypes}
	}

	if len(sockets) > 0 {
		filter["main.socket"] = bson.M{"$in": sockets}
	}

	if len(pcie) > 0 {
		filter["pci-e"] = bson.M{"$in": pcie}
	}

	logger.Infof("Set filter: %v", filter)
	return filter
}

func getProductsFilter(productType string, r *http.Request) bson.M {
	var filter bson.M

	switch productType {
	case "cpu":
		filter = filterCpu(r)
	default:
		filter = bson.M{}
	}
	return filter
}

func getBuildFromCookie(r *http.Request) ([]data.ProductHeader, error) {
	logger := logging.GetLogger()
	userID := data.ShowUser(r).ID.Hex()
	buildCookie, err := r.Cookie("build" + userID)

	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return []data.ProductHeader{}, nil
		}
		return nil, err
	}

	decodedValue, err := url.QueryUnescape(buildCookie.Value)
	if err != nil {
		return nil, err
	}

	fmt.Println("Decoded Build Cookie Value:", decodedValue)

	var build []data.ProductHeader
	err = json.Unmarshal([]byte(decodedValue), &build)
	if err != nil {
		logger.Errorf("Error unmarshalling build cookie: %v", err)
		logger.Errorf("Build Cookie Value (Raw): %v", decodedValue)

		return []data.ProductHeader{}, nil
	}

	return build, nil
}

func saveBuildToCookie(w http.ResponseWriter, r *http.Request, build []data.ProductHeader) error {
	logger := logging.GetLogger()
	buildJSON, err := json.Marshal(build)
	if err != nil {
		return err
	}

	userID := data.ShowUser(r).ID.Hex()

	encodedValue := url.QueryEscape(string(buildJSON))

	http.SetCookie(w, &http.Cookie{
		Name:    "build" + userID,
		Value:   encodedValue,
		Expires: time.Now().Add(240 * time.Hour),
		Path:    "/",
	})

	logger.Infof("Cookie with name: %s was saved.", userID)
	return nil
}

/*func GetBuild(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	userBuild, err := getCartFromCookie(r)
	if err != nil {
		HandleError(errors.New("error getting build"), logger, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(userBuild); err != nil {
		HandleError(errors.New("error encoding build data"), logger, w)
		return
	}
}*/

func AddToBuild(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	dbName := "shop"
	productType := r.FormValue("productType")

	collectionName, product, err := defineStruct(productType)
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

	result, err := data.GetProductById(dbName, collectionName, ObjID)
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

	item, err := extractProductHeaderFromProduct(productType, product)
	if err != nil {
		logger.Errorf("Error extracting item from product: %v", err)
		HandleError(err, logger, w)
		return
	}

	var userBuild []data.ProductHeader
	userBuild, err = getBuildFromCookie(r)
	if err != nil {
		logger.Infof("Error getting build: %v", err)
		HandleError(err, logger, w)
		return
	}

	found := false
	for i, buildItem := range userBuild {
		if buildItem.ID == item.ID {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
			logger.Infof("Product is already in build: %v", item)
			fmt.Println(userBuild)
			return
		}
		if buildItem.ProductType == item.ProductType {
			userBuild[i] = item
			found = true
			logger.Infof("Item in build with type %s was replaced by: %v", item.ProductType, item)
			break
		}
	}

	if !found {
		userBuild = append(userBuild, item)
	}

	err = saveBuildToCookie(w, r, userBuild)
	if err != nil {
		logger.Infof("Error saving build to cookie: %v", err)
		http.Error(w, "Error saving build", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	logger.Infof("Product added to build: %v", item)
	fmt.Println(userBuild)
}

func extractProductHeaderFromProduct(productType string, product interface{}) (data.ProductHeader, error) {
	v := reflect.ValueOf(product)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	idField := v.FieldByName("ID")

	if !idField.IsValid() {
		return data.ProductHeader{}, errors.New("invalid product structure")
	}

	return data.ProductHeader{
		ID:          idField.Interface().(primitive.ObjectID),
		ProductType: productType,
	}, nil
}

func DeleteFromBuild(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productType := r.FormValue("productType")
	userBuild, err := getBuildFromCookie(r)
	if err != nil {
		HandleError(errors.New("error getting build"), logger, w)
	}
	fmt.Println(productType)

	var updatedBuild []data.ProductHeader
	for _, item := range userBuild {
		if item.ProductType != productType {
			updatedBuild = append(updatedBuild, item)
		}
	}

	err = saveBuildToCookie(w, r, updatedBuild)
	if err != nil {
		logger.Errorf("Error saving updated build to cookie: %v", err)
		HandleError(errors.New("error saving updated build"), logger, w)
		return
	}

	logger.Infof("Product deleted from cart: %v", productType)
}

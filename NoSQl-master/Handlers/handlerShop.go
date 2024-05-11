package Handlers

import (
	"MongoDb/internal/data"
	"MongoDb/pkg/logging"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/http"
	"strconv"
)

//TODO 1 closing connections (1/2 DONE)
//TODO 2 make single function for products list (DONE)
//TODO make single delete, modify functions

func Shop(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/shop.html"))
	tmpl.Execute(w, data.ShowUser())
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

		dataToSend := []interface{}{items, data.ShowUser().Name}

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

		dataToSend := []interface{}{items, data.ShowUser().Name}

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

		dataToSend := []interface{}{items, data.ShowUser().Name}

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())

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
	tmpl.Execute(w, data.ShowUser())
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
	value := r.URL.Query().Get("element")
	var tmplName string
	var itemsStruct []interface{}
	dbName := "shop"
	var collectionName string

	switch value {
	case "cpu":
		collectionName = "cpu"
		tmplName = "listCpu.html"
	case "cooling":
		collectionName = "cooling"
		tmplName = "listCooling.html"
	case "motherboard":
		collectionName = "motherboard"
		tmplName = "listMotherboard.html"
	case "housing":
		collectionName = "housing"
		tmplName = "listHousing.html"
	case "hdd":
		collectionName = "hdd"
		tmplName = "listHdd.html"
	case "ssd":
		collectionName = "ssd"
		tmplName = "listSsd.html"
	case "powersupply":
		collectionName = "powersupply"
		tmplName = "listPowersupply.html"
	case "gpu":
		collectionName = "gpu"
		tmplName = "listGpu.html"
	case "ram":
		collectionName = "ram"
		tmplName = "listRam.html"
	default:
		logger.Infof("Error: failed to List Products!")
		return
	}

	err := data.Init(dbName, collectionName)
	if err != nil {
		return
	}
	defer data.CloseConnection()

	tmpl := template.Must(template.ParseFiles("html/" + tmplName))

	filter := bson.M{}

	cur, err := data.Collection.Find(context.TODO(), filter)
	if err != nil {
		logger.Infof("error: %v", err)
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		switch value {
		case "cpu":
			var cpuItem data.Cpu
			err := cur.Decode(&cpuItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			itemsStruct = append(itemsStruct, cpuItem)
		case "cooling":
			var coolingItem data.Cooling
			err := cur.Decode(&coolingItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			itemsStruct = append(itemsStruct, coolingItem)
		case "motherboard":
			var motherboardItem data.Motherboard
			err := cur.Decode(&motherboardItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			itemsStruct = append(itemsStruct, motherboardItem)
		case "housing":
			var housingItem data.Housing
			err := cur.Decode(&housingItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			itemsStruct = append(itemsStruct, housingItem)
		case "hdd":
			var hddItem data.Hdd
			err := cur.Decode(&hddItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			itemsStruct = append(itemsStruct, hddItem)
		case "ssd":
			var ssdItem data.Ssd
			err := cur.Decode(&ssdItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			itemsStruct = append(itemsStruct, ssdItem)
		case "powersupply":
			var powerSupplyItem data.PowerSupply
			err := cur.Decode(&powerSupplyItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			itemsStruct = append(itemsStruct, powerSupplyItem)
		case "gpu":
			var gpuItem data.Gpu
			err := cur.Decode(&gpuItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			itemsStruct = append(itemsStruct, gpuItem)
		case "ram":
			var ramItem data.Ram
			err := cur.Decode(&ramItem)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			itemsStruct = append(itemsStruct, ramItem)
		}
	}

	if err := cur.Err(); err != nil {
		logger.Infof("error: %v", err)
	}

	logger.Infof("Found multiple items: %v", len(itemsStruct))

	dataToSend := []interface{}{itemsStruct, data.ShowUser().Name}

	err = tmpl.Execute(w, dataToSend)
	if err != nil {
		logger.Infof("error: %v", err)
	}
}

func ListProductInfo(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	value := r.FormValue("showProduct")[0:3]
	var tmplName string
	dbName := "shop"
	var collectionName string
	var item interface{}

	switch value {
	case "cpu":
		collectionName = "cpu"
		tmplName = "fullListCpu.html"
		item = data.Cpu{}
	case "clg":
		collectionName = "cooling"
		tmplName = "fullListCooling.html"
		item = data.Cooling{}
	case "mbd":
		collectionName = "motherboard"
		tmplName = "fullListMotherboard.html"
		item = data.Motherboard{}
	case "hsg":
		collectionName = "housing"
		tmplName = "fullListHousing.html"
		item = data.Housing{}
	case "hdd":
		collectionName = "hdd"
		tmplName = "fullListHdd.html"
		item = data.Hdd{}
	case "ssd":
		collectionName = "ssd"
		tmplName = "fullListSsd.html"
		item = data.Ssd{}
	case "pwr":
		collectionName = "powersupply"
		tmplName = "fullListPowersupply.html"
		item = data.PowerSupply{}
	case "gpu":
		collectionName = "gpu"
		tmplName = "fullListGpu.html"
		item = data.Gpu{}
	case "ram":
		collectionName = "ram"
		tmplName = "fullListRam.html"
		item = data.Ram{}
	default:
		logger.Infof("Error: failed to Show Product Info!")
		return
	}

	err := data.Init(dbName, collectionName)
	if err != nil {
		return
	}
	defer data.CloseConnection()

	ObjID, err := primitive.ObjectIDFromHex(r.FormValue("showProduct")[13:37])
	filter := bson.M{"_id": ObjID}

	cur, err := data.Collection.Find(context.TODO(), filter)
	if err != nil {
		logger.Infof("error: %v", err)
	}
	defer cur.Close(context.TODO())

	tmpl := template.Must(template.ParseFiles("html/" + tmplName))

	for cur.Next(context.TODO()) {
		switch value {
		case "cpu":
			var itemCpu data.Cpu
			err := cur.Decode(&itemCpu)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			item = itemCpu
		case "clg":
			var itemCooling data.Cooling
			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			item = itemCooling
		case "mbd":
			var itemMotherboard data.Motherboard
			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			item = itemMotherboard
		case "hsg":
			var itemHosuing data.Housing
			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			item = itemHosuing
		case "hdd":
			var itemHdd data.Hdd
			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			item = itemHdd
		case "ssd":
			var itemSsd data.Ssd
			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			item = itemSsd
		case "pwr":
			var itemPowerSupply data.PowerSupply
			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			item = itemPowerSupply
		case "gpu":
			var itemGpu data.Gpu
			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			item = itemGpu
		case "ram":
			var itemRam data.Ram
			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error: %v", err)
				continue
			}
			item = itemRam
		}
	}

	if err := cur.Err(); err != nil {
		logger.Infof("error: %v", err)
	}

	logger.Infof("Found Item: %v", item)

	err = tmpl.Execute(w, item)
	if err != nil {
		logger.Infof("error: %v", err)
	}
}

func FilterCpu(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	r.ParseForm()

	var manufacturers []string
	var categories []string
	var cores []int
	var priceFrom, priceTo float64
	var ramTypes []string
	var sockets []string
	var pcieVersions []float64

	manufacturers = r.Form["manufacturer"]
	categories = r.Form["category"]
	coreValues := r.Form["cores"]
	for _, value := range coreValues {
		core, err := strconv.Atoi(value)
		if err != nil {
			logger.Infof("Error parsing core value: %v", err)
			continue
		}
		cores = append(cores, core)
	}
	priceFromStr := r.Form.Get("price-from")
	priceToStr := r.Form.Get("price-to")

	var err error
	if priceFromStr != "" {
		priceFrom, err = strconv.ParseFloat(priceFromStr, 64)
		if err != nil {
			logger.Infof("Error parsing price from: %v", err)
		}
	} else {
		priceFrom = 0
	}

	if priceToStr != "" {
		priceTo, err = strconv.ParseFloat(priceToStr, 64)
		if err != nil {
			logger.Infof("Error parsing price to: %v", err)
		}
	} else {
		priceTo = 999999
	}

	ramTypes = r.Form["ram-type"]
	sockets = r.Form["socket"]
	pcieValues := r.Form["pcie-version"]
	for _, value := range pcieValues {
		pcie, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.Infof("Error parsing PCIE version: %v", err)
			continue
		}
		pcieVersions = append(pcieVersions, pcie)
	}

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
		filter["manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(categories) > 0 {
		categoryFilter := bson.M{"$in": categories}
		filter["main.category"] = categoryFilter
	}

	if len(cores) > 0 {
		filter["cores.p-cores"] = bson.M{"$in": cores}
	}

	if priceFrom != 0 || priceTo != 0 {
		filter["price"] = bson.M{"$gte": priceFrom, "$lte": priceTo}
	}

	if len(ramTypes) > 0 {
		filter["ram.type"] = bson.M{"$in": ramTypes}
	}

	if len(sockets) > 0 {
		filter["main.socket"] = bson.M{"$in": sockets}
	}

	if len(pcieVersions) > 0 {
		filter["pci-e"] = bson.M{"$in": pcieVersions}
	}

	logger.Infof("Set filter: %v", filter)

	// Connect to MongoDB and fetch filtered CPU data
	dbName := "shop"
	collectionName := "cpu"
	err = data.Init(dbName, collectionName)
	if err != nil {
		return
	}
	defer data.CloseConnection()

	cur, err := data.Collection.Find(context.Background(), filter)
	if err != nil {
		logger.Infof("Error finding CPUs: %v", err)
		return
	}
	defer cur.Close(context.Background())

	var cpuItems []data.Cpu
	for cur.Next(context.Background()) {
		var cpuItem data.Cpu
		err := cur.Decode(&cpuItem)
		if err != nil {
			logger.Infof("Error decoding CPU item: %v", err)
			continue
		}
		cpuItems = append(cpuItems, cpuItem)
	}
	if err := cur.Err(); err != nil {
		logger.Infof("Error iterating cursor: %v", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/listCpu.html"))
	logger.Infof("Found multiple items: %v", len(cpuItems))

	dataToSend := []interface{}{cpuItems, data.ShowUser().Name}

	err = tmpl.Execute(w, dataToSend)
	if err != nil {
		logger.Infof("error: %v", err)
	}

	if err != nil {
		logger.Infof("Error executing template: %v", err)
		return
	}
}

func ShowProfile(w http.ResponseWriter, r *http.Request) {
	//logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/userProfile.html"))
	tmpl.Execute(w, data.ShowUser())
}

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

var ItemsCpu []data.Cpu
var ItemCpu data.Cpu

/*Shop function*/
func Shop(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/shop.html"))
	tmpl.Execute(w, data.GetUser(w))
}

/*
func listCpu() (itmArr interface{}, itm interface{}) {
	var itemsCpu []data.Cpu
	var itemCpu data.Cpu
	return itemsCpu, itemCpu
}
*/

func ComparisonCpuMb(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listMotherboard.html"))
	tmpl.Execute(w, data.GetUser(w))

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
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonCpuRam(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listRam.html"))
	tmpl.Execute(w, data.GetUser(w))

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

func ComparisonCpuCooling(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listCooling.html"))
	tmpl.Execute(w, data.GetUser(w))

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
		fmt.Println("Found multiple items:", len(items))

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ComparisonMbCpu(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/listCpu.html"))
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))

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
	tmpl.Execute(w, data.GetUser(w))
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

		main := data.MainCpu{
			Category:   r.FormValue("category"),
			Model:      r.FormValue("model"),
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
			Manufacturer:   r.FormValue("man"),
			Main:           main,
			Cores:          cores,
			ClockFrequency: clockFrequency,
			Ram:            ram,
			Tdp:            tdp,
			Graphics:       r.FormValue("graphics"),
			PciE:           pcie,
			MaxTemperature: maxTemp,
			Price:          price,
			Discount:       discount,
			Amount:         amount,
		}

		_, err := data.Collection.InsertOne(context.TODO(), recordCpu)
		if err != nil {
			logger.Infof("A bulk write error occurred: %v", err)
		} else {
			logger.Infof("CPU record with ID: %s was CREATED!", recordCpu.ID)
		}

	}
}

func ModifyCpuForm(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/modifyCpu.html"))

	data.Init("shop", "cpu")
	var items []data.Cpu

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("modify")[10:34])
		filter := bson.M{"_id": ObjID}

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

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

func ModifyCpu(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	data.Init("shop", "cpu")

	if r.Method == http.MethodPost {
		year, _ := strconv.Atoi(r.FormValue("year"))
		pcores, _ := strconv.Atoi(r.FormValue("pcores"))
		ecores, _ := strconv.Atoi(r.FormValue("ecores"))
		threads, _ := strconv.Atoi(r.FormValue("threads"))
		techPr, _ := strconv.Atoi(r.FormValue("techPr"))
		pcoresBase, _ := strconv.ParseFloat(r.FormValue("pcoresbase"), 64)
		pcoresBoost, _ := strconv.ParseFloat(r.FormValue("pcoresturbo"), 64)
		ecoresBase, _ := strconv.ParseFloat(r.FormValue("ecoresbase"), 64)
		ecoresBoost, _ := strconv.ParseFloat(r.FormValue("ecoresturbo"), 64)
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

		main := data.MainCpu{
			Category:   r.FormValue("category"),
			Model:      r.FormValue("model"),
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
			Manufacturer:   r.FormValue("man"),
			Main:           main,
			Cores:          cores,
			ClockFrequency: clockFrequency,
			Ram:            ram,
			Tdp:            tdp,
			Graphics:       r.FormValue("graphics"),
			PciE:           pcie,
			MaxTemperature: maxTemp,
			Price:          price,
			Discount:       discount,
			Amount:         amount,
		}

		update := bson.M{"$set": bson.M{
			"manufacturer":    recordCpu.Manufacturer,
			"main":            recordCpu.Main,
			"cores":           recordCpu.Cores,
			"clock_frequency": recordCpu.ClockFrequency,
			"ram":             recordCpu.Ram,
			"tdp":             recordCpu.Tdp,
			"graphics":        recordCpu.Graphics,
			"pci-e":           recordCpu.PciE,
			"max_temperature": recordCpu.MaxTemperature,
			"price":           recordCpu.Price,
			"discount":        recordCpu.Discount,
			"amount":          recordCpu.Amount,
		}}

		_, err = data.Collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			logger.Infof("A bulk write error occurred: %v", err)
		} else {
			logger.Infof("CPU record with ID: %s was UPDATED!", ObjID)
		}

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
}

/*func listObjects(w http.ResponseWriter, r *http.Request, value string) {
	logger := logging.GetLogger()

	tmpl := template.Must(template.ParseFiles("html/listCpu.html"))
	tmpl.Execute(w, data.GetUser(w))

	data.Init("shop", "cpu")
	var items interface{}
	var item interface{}

	if r.Method == http.MethodPost {

		filter := bson.M{}

		cur, err := data.Collection.Find(context.Background(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.Background())

		for cur.Next(context.Background()) {
			// convert item to a data.Cpu if necessary
			if item == nil {
				item = data.Cpu{}
			}
			cpu, ok := item.(data.Cpu)
			if !ok {
				logger.Infof("item is not of type data.Cpu")
				return
			}

			err := cur.Decode(&cpu)
			if err != nil {
				logger.Infof("error :", err)
			}

			// convert items to a slice of data.Cpu if necessary
			if items == nil {
				items = []data.Cpu{}
			}
			slice, ok := items.([]data.Cpu)
			if !ok {
				logger.Infof("items is not of type []data.Cpu")
				return
			}

			slice = append(slice, cpu)
			items = slice
		}

		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}

		data.Init("test", "users")
	}
}*/

/*Full List of Gpu*/
func ListGpu(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/fullListGpu.html"))
	tmpl.Execute(w, data.GetUser(w))

	data.Init("shop", "gpu")
	var items []data.Gpu

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("ssdbtn")[10:34])
		filter := bson.M{"_id": ObjID}

		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.Gpu

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", items)

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

/*Full List of PowerSupply*/
func ListPowerSupply(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/fullListPowerSupply.html"))
	tmpl.Execute(w, data.GetUser(w))

	data.Init("shop", "powersupply")
	var items []data.PowerSupply

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("ssdbtn")[10:34])
		filter := bson.M{"_id": ObjID}

		cur, err := data.Collection.Find(context.TODO(), filter)
		if err != nil {
			logger.Infof("error:", err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var item data.PowerSupply

			err := cur.Decode(&item)
			if err != nil {
				logger.Infof("error :", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			logger.Infof("error :", err)
		}
		fmt.Println("Found multiple items:", items)

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

/*Full List of Motherboard*/
func ListMotherboard(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/fullListMotherboard.html"))
	tmpl.Execute(w, data.GetUser(w))

	data.Init("shop", "motherboard")
	var items []data.Motherboard

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("ssdbtn")[10:34])
		filter := bson.M{"_id": ObjID}

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
		fmt.Println("Found multiple items:", items)

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

/*Full List of Ram*/
func ListRam(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/fullListRam.html"))
	tmpl.Execute(w, data.GetUser(w))

	data.Init("shop", "ram")
	var items []data.Ram

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("ssdbtn")[10:34])
		filter := bson.M{"_id": ObjID}

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
		fmt.Println("Found multiple items:", items)

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

/*Full List of Housing*/
func ListHousing(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/fullListHousing.html"))
	tmpl.Execute(w, data.GetUser(w))

	data.Init("shop", "housing")
	var items []data.Housing

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("ssdbtn")[10:34])
		filter := bson.M{"_id": ObjID}

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
		fmt.Println("Found multiple items:", items)

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

/*Full List of Hdd*/
func ListHdd(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/fullListHdd.html"))
	tmpl.Execute(w, data.GetUser(w))

	data.Init("shop", "hdd")
	var items []data.Hdd

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("ssdbtn")[10:34])
		filter := bson.M{"_id": ObjID}

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
		fmt.Println("Found multiple items:", items)

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

/*Full List of Cooling*/
func ListCooling(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/fullListCooling.html"))
	tmpl.Execute(w, data.GetUser(w))

	data.Init("shop", "cooling")
	var items []data.Cooling

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("ssdbtn")[10:34])
		filter := bson.M{"_id": ObjID}

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
		fmt.Println("Found multiple items:", items)

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

/*Full List of Cpu*/
func ListCpu(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/fullListCpu.html"))

	data.Init("shop", "cpu")
	var items []data.Cpu

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("ssdbtn")[10:34])
		filter := bson.M{"_id": ObjID}

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

		err = tmpl.Execute(w, items)
		if err != nil {
			logger.Infof("error :", err)
		}
		data.Init("test", "users")
	}
}

/*Full List of Ssd*/
func ListSsd(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/fullListSsd.html"))
	tmpl.Execute(w, data.GetUser(w))

	data.Init("shop", "ssd")
	var items []data.Ssd

	if r.Method == http.MethodPost {

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("ssdbtn")[10:34])
		filter := bson.M{"_id": ObjID}

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

/*Function to list elements that we have*/
func List(w http.ResponseWriter, r *http.Request) {

	logger := logging.GetLogger()

	if r.Method == http.MethodPost {
		value := r.FormValue("element")

		switch value {
		case "cpu":
			tmpl := template.Must(template.ParseFiles("html/listCpu.html"))
			data.Init("shop", "cpu")
			var items []data.Cpu

			filter := bson.M{}

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
		case "cooling":
			tmpl := template.Must(template.ParseFiles("html/listCooling.html"))
			data.Init("shop", "cooling")
			var items []data.Cooling

			filter := bson.M{}

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
			fmt.Println("Found multiple items:", len(items))

			err = tmpl.Execute(w, items)
			if err != nil {
				logger.Infof("error :", err)
			}
			data.Init("test", "users")
		case "hdd":
			tmpl := template.Must(template.ParseFiles("html/listHdd.html"))
			tmpl.Execute(w, data.GetUser(w))
			data.Init("shop", "hdd")
			var items []data.Hdd

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

		case "housing":
			tmpl := template.Must(template.ParseFiles("html/listHousing.html"))
			tmpl.Execute(w, data.GetUser(w))
			data.Init("shop", "housing")
			var items []data.Housing

			filter := bson.M{}

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
		case "ram":
			tmpl := template.Must(template.ParseFiles("html/listRam.html"))
			tmpl.Execute(w, data.GetUser(w))
			data.Init("shop", "ram")
			var items []data.Ram

			filter := bson.M{}

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

		case "motherboard":
			tmpl := template.Must(template.ParseFiles("html/listMotherboard.html"))
			tmpl.Execute(w, data.GetUser(w))
			data.Init("shop", "motherboard")
			var items []data.Motherboard

			filter := bson.M{}

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

		case "ssd":
			tmpl := template.Must(template.ParseFiles("html/listSsd.html"))
			tmpl.Execute(w, data.GetUser(w))
			data.Init("shop", "ssd")
			var items []data.Ssd

			filter := bson.M{}

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

		case "powersupply":
			tmpl := template.Must(template.ParseFiles("html/listPowersupply.html"))
			tmpl.Execute(w, data.GetUser(w))
			data.Init("shop", "powersupply")
			var items []data.PowerSupply

			filter := bson.M{}

			cur, err := data.Collection.Find(context.TODO(), filter)
			if err != nil {
				logger.Infof("error:", err)
			}
			defer cur.Close(context.TODO())

			for cur.Next(context.TODO()) {
				var item data.PowerSupply

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

		case "gpu":
			tmpl := template.Must(template.ParseFiles("html/listGpu.html"))
			tmpl.Execute(w, data.GetUser(w))
			data.Init("shop", "gpu")
			var items []data.Gpu

			filter := bson.M{}

			cur, err := data.Collection.Find(context.TODO(), filter)
			if err != nil {
				logger.Infof("error:", err)
			}
			defer cur.Close(context.TODO())

			for cur.Next(context.TODO()) {
				var item data.Gpu

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
}

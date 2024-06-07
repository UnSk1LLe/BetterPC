package Handlers

import (
	"MongoDb/internal/data"
	"MongoDb/internal/filters"
	"MongoDb/pkg/logging"
	"MongoDb/pkg/session"
	"MongoDb/pkg/templateFunctions"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

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

func AddNewProductChoice(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/addProductChoice.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func AddNewProductForm(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productType := r.FormValue("productType")
	var tmpl *template.Template
	switch productType {
	case "cpu":
		tmpl = template.Must(template.ParseFiles("html/addCpu.html"))
	default:
		logger.Errorf("invalid productType: %s", productType)
		HandleError(fmt.Errorf("invalid productType: %s", productType), logger, w)
		return
	}
	err := tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func AddNewProduct(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productType := r.FormValue("productType")

	switch productType {
	case "cpu":
		addCpu(w, r)
	default:
		logger.Errorf("invalid productType: %s", productType)
		HandleError(fmt.Errorf("invalid productType: %s", productType), logger, w)
		return
	}
	_ = showMessage("/shop", "New product created!", w)
}

func addCpu(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ddr4MaxFr := 0
		ddr5MaxFr := 0

		types := r.Form["type"]

		for _, t := range types {
			if t == "DDR4" {
				ddr4MaxFr, _ = strconv.Atoi(r.FormValue("ddr4MaxFr"))
			}
			if t == "DDR5" {
				ddr5MaxFr, _ = strconv.Atoi(r.FormValue("ddr5MaxFr"))
			}
		}

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
			MaxFrequency: []int{ddr4MaxFr, ddr5MaxFr},
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

		_, err := data.CpuCollection.InsertOne(context.TODO(), recordCpu)
		if err != nil {
			logger.Infof("A bulk write error occurred: %v", err)
		} else {
			logger.Infof("CPU record with ID: %s was CREATED!", recordCpu.ID)
		}
	}
}

func ModifyProductForm(w http.ResponseWriter, r *http.Request) {
	var tmpl *template.Template
	logger := logging.GetLogger()
	productType := r.FormValue("productType")
	productID := r.FormValue("productID")

	ObjID, err := primitive.ObjectIDFromHex(productID[10:34])
	if err != nil {
		logger.Errorf("invalid productID: %s", productID)
		HandleError(fmt.Errorf("invalid productID: %s", productID), logger, w)
		return
	}

	_, product, err := defineStruct(productType)
	if err != nil {
		logger.Errorf("invalid productType: %s", productType)
		HandleError(fmt.Errorf("invalid productType: %s", productType), logger, w)
		return
	}

	switch productType {
	case "cpu":
		tmpl = template.Must(template.ParseFiles("html/modifyCpu.html"))
		var item data.Cpu
		_, err = getAndDecodeProduct(data.CpuCollection, &item, ObjID)
		product = item
	default:
		_ = showMessage("/shop", "404 Not found!", w)
		return
	}
	if err != nil {
		logger.Errorf("Error finding product of type: %s", productType)
		HandleError(fmt.Errorf("error finding product of type: %s", productType), logger, w)
		return
	}

	dataToSend := struct {
		ProductType string
		Product     interface{}
	}{
		ProductType: productType,
		Product:     product,
	}

	err = tmpl.Execute(w, dataToSend)
	if err != nil {
		return
	}
}

func ModifyProduct(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	productType := r.FormValue("productType")
	ObjID, err := primitive.ObjectIDFromHex(r.FormValue("productID")[10:34])
	if err != nil {
		logger.Errorf("invalid productID: %s", r.FormValue("modifyCpu"))
		HandleError(fmt.Errorf("invalid productID: %s", r.FormValue("modifyCpu")), logger, w)
		return
	}

	switch productType {
	case "cpu":
		err = modifyCpu(ObjID, r)
	default:
		_ = showMessage("/shop", "Internal server error!", w)
	}
	if err != nil {
		logger.Errorf("Error updating product <%s> with ID", productType)
		HandleError(err, logger, w)
		return
	}
	_ = showMessage("/shop", "Product updated!", w)
	return
}

func modifyCpu(ID primitive.ObjectID, r *http.Request) error {
	logger := logging.GetLogger()

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return err
		}

		ddr4MaxFr := 0
		ddr5MaxFr := 0

		types := r.Form["type"]

		for _, t := range types {
			if t == "DDR4" {
				ddr4MaxFr, _ = strconv.Atoi(r.FormValue("ddr4MaxFr"))
			}
			if t == "DDR5" {
				ddr5MaxFr, _ = strconv.Atoi(r.FormValue("ddr5MaxFr"))
			}
		}

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
			MaxFrequency: []int{ddr4MaxFr, ddr5MaxFr},
			MaxCapacity:  ramMaxCap,
		}

		filter := bson.M{"_id": ID}

		recordCpu := data.Cpu{
			ID:             ID,
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

		fmt.Printf("CPU %v", recordCpu)
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

		res, err := data.UpdateProduct(data.CpuCollection, filter, update)
		fmt.Println(res)
		if err != nil {
			logger.Errorf("A bulk write error occurred: %v", err)
			return err
		}
	} else {
		return errors.New("not a post method")
	}
	logger.Infof("CPU record with ID: %s was UPDATED!", ID)
	return nil
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productType := r.FormValue("productType")
	productID := r.FormValue("productID")
	ObjID, err := primitive.ObjectIDFromHex(productID[10:34])
	if err != nil {
		logger.Errorf("Error retreving objectID from hex: %v", err)
		HandleError(fmt.Errorf("error retreving objectID from hex: %v", err), logger, w)
		return
	}
	_, err = data.DeleteProductById(productType, ObjID)
	if err != nil {
		HandleError(err, logger, w)
	}
	logger.Infof("Product deleted successfully")
	message := "Product <" + productType + "> with ID: " + productID + " DELETED successfully!"
	_ = showMessage("/shop", message, w)
}

func getCompatibilityFilter(productType string, r *http.Request) (bson.M, error) {
	logger := logging.GetLogger()
	build, err := getFullBuild(r)
	filter := bson.M{}
	var conditions []bson.M
	if err != nil {
		logger.Errorf("A bulk write error occurred: %v", err)
		return bson.M{}, err
	}

	switch productType {
	case "cpu":
		if !data.IsZero(build.Motherboard) {
			conditions = append(conditions, bson.M{"main.socket": build.Motherboard.Socket})
		}
		if !data.IsZero(build.RAM) {
			if build.RAM.Type == "DDR4" {
				conditions = append(conditions, bson.M{"ram.max_frequency.0": bson.M{"$gte": build.RAM.Frequency}})
			} else if build.RAM.Type == "DDR5" {
				conditions = append(conditions, bson.M{"ram.max_frequency.1": bson.M{"$gte": build.RAM.Frequency}})
			}
			conditions = append(conditions, bson.M{"ram.max_capacity": bson.M{"$gte": build.RAM.Capacity * build.RAM.Number}})
		}
		if !data.IsZero(build.Cooling) {
			conditions = append(conditions, bson.M{"main.socket": bson.M{"$in": build.Cooling.Sockets}, "tdp": bson.M{"$lte": build.Cooling.Tdp}})
		}
	case "motherboard":
		if !data.IsZero(build.CPU) {
			conditions = append(conditions, bson.M{"socket": build.CPU.Main.Socket})
		}
		if !data.IsZero(build.Housing) {
			conditions = append(conditions, bson.M{"form_factor": build.Housing.MbFormFactor})
		}
		if !data.IsZero(build.RAM) {
			conditions = append(conditions, bson.M{"ram.type": build.RAM.Type, "ram.max_frequency": bson.M{"$gte": build.RAM.Frequency},
				"ram.max_capacity": bson.M{"$gte": build.RAM.Capacity * build.RAM.Number}, "ram.slots": bson.M{"$gte": build.RAM.Number}})
		}
	case "ram":
		if !data.IsZero(build.Motherboard) {
			conditions = append(conditions, bson.M{"type": build.Motherboard.Ram.Type, "frequency": bson.M{"$lte": build.Motherboard.Ram.MaxFrequency}, "form_factor": "DIMM"})
		}
		if !data.IsZero(build.CPU) {
			if build.CPU.Ram.MaxFrequency[0] != 0 {
				conditions = append(conditions, bson.M{"frequency": bson.M{"$lte": build.CPU.Ram.MaxFrequency[0]}})
			}
			if build.CPU.Ram.MaxFrequency[1] != 0 {
				conditions = append(conditions, bson.M{"frequency": bson.M{"$lte": build.CPU.Ram.MaxFrequency[1]}})
			}
			conditions = append(conditions, bson.M{"$expr": bson.M{"$lte": bson.A{bson.M{"$multiply": bson.A{"$capacity", "$number"}}, build.CPU.Ram.MaxCapacity}}})
		}
	case "gpu":
		if !data.IsZero(build.PowerSupply) {
			conditions = append(conditions, bson.M{"tdp_r": bson.M{"$lte": build.PowerSupply.OutputPower}})
		}
		if !data.IsZero(build.Housing) {
			conditions = append(conditions, bson.M{"size.0": bson.M{"$lte": build.Housing.GraphicCardSize}})
		}
	case "ssd":
		if !data.IsZero(build.Motherboard) {
			if build.Motherboard.Interfaces.M2 > 0 && build.Motherboard.Interfaces.Sata3 > 0 {
				filter = bson.M{"$or": []bson.M{{"form_factor": bson.M{"$regex": "M.2"}}, {"interface": "SATA 3"}}}
			} else if build.Motherboard.Interfaces.M2 > 0 {
				conditions = append(conditions, bson.M{"form_factor": bson.M{"$regex": "M.2"}})
			} else if build.Motherboard.Interfaces.Sata3 > 0 {
				conditions = append(conditions, bson.M{"interface": "SATA 3"})
			}
		}
	case "hdd":
		if !data.IsZero(build.Motherboard) {
			if build.Motherboard.Interfaces.Sata3 > 0 {
				filter = bson.M{"interface": "SATA 3"}
			}
		}
	case "cooling":
		if !data.IsZero(build.CPU) {
			conditions = append(conditions, bson.M{"sockets": build.CPU.Main.Socket})
			conditions = append(conditions, bson.M{"tdp": bson.M{"$gte": build.CPU.Tdp}})
		}
		if !data.IsZero(build.Motherboard) {
			conditions = append(conditions, bson.M{"sockets": build.Motherboard.Socket})
		}
		if !data.IsZero(build.Housing) {
			conditions = append(conditions, bson.M{"height": bson.M{"$lte": build.Housing.CoolerHeight}})
		}
	case "powersupply":
		if !data.IsZero(build.GPU) {
			conditions = append(conditions, bson.M{"output_power": bson.M{"$gte": build.GPU.TdpR}})
		}
		if !data.IsZero(build.Housing) {
			conditions = append(conditions, bson.M{"ps_form_factor": build.Housing.PsFormFactor})
		}
	case "housing":
		if !data.IsZero(build.Motherboard) {
			conditions = append(conditions, bson.M{"mb_form_factor": build.Motherboard.FormFactor})
		}
		if !data.IsZero(build.PowerSupply) {
			conditions = append(conditions, bson.M{"ps_form_factor": build.PowerSupply.FormFactor})
		}
		if !data.IsZero(build.SSD) {
			if build.SSD.FormFactor == "2.5" {
				conditions = append(conditions, bson.M{"drive_bays.2_5": bson.M{"$gte": 0}})
			}
		}
		if !data.IsZero(build.HDD) {
			conditions = append(conditions, bson.M{"drive_bays.3_5": bson.M{"$gte": 0}})
		}
	default:
		logger.Errorf("Unknown product type: %s", productType)
	}

	if len(conditions) > 0 {
		filter["$and"] = conditions
	}

	fmt.Println(filter)
	return filter, nil
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productType := r.URL.Query().Get("productType")
	listCompatibleOnly := r.URL.Query().Get("listCompatibleOnly")
	searchQuery := r.URL.Query().Get("search")
	fmt.Println(1, searchQuery)
	var buildFilter bson.M
	var err error
	search := false
	pcBuilder := false

	if listCompatibleOnly == "true" {
		pcBuilder = true
		buildFilter, err = getCompatibilityFilter(productType, r)
		if err != nil {
			buildFilter = bson.M{}
		}
	} else {
		buildFilter = bson.M{}
	}

	pageNumber, err := strconv.Atoi(r.FormValue("pageNumber"))
	if err != nil {
		pageNumber = 1
		logger.Infof("Failed to parse page number: %s, so it set to 1", err)
	}

	var productsList []data.Product

	tmpl := template.New("listProducts.html").Funcs(templateFunctions.TmplFuncs)
	tmpl, err = tmpl.ParseFiles("html/listProducts.html")
	if err != nil {
		logger.Errorf("Failed to parse template: %s", err)
		return
	}

	//tmpl := template.Must(template.ParseFiles("html/listProducts.html"))

	filter := getProductsFilter(productType, r)

	searchFilter := filters.SearchProducts(searchQuery)
	if searchFilter != nil {
		search = true
	}
	finalFilter := bson.M{"$and": []bson.M{buildFilter, filter, searchFilter}}

	perPageLimit := 7
	productsList, productsNumber, err := productsListing(productType, finalFilter, pageNumber, perPageLimit, w)
	if err != nil {
		HandleError(err, logger, w)
		return
	}

	totalPages := int(math.Ceil(float64(productsNumber) / float64(perPageLimit)))

	user, _ := data.GetUserBySessionToken(session.GetSessionTokenFromCookie(r))
	build, err := getBuild(r)

	if !search {
		searchQuery = ""
	}

	dataToSend := struct {
		ProductType  string
		ProductsList []data.Product
		Build        *data.Build
		User         data.User
		PcBuilder    bool
		CurrentPage  int
		TotalPages   int
		SearchQuery  string
	}{
		ProductType:  productType,
		ProductsList: productsList,
		Build:        build,
		User:         user,
		PcBuilder:    pcBuilder,
		CurrentPage:  pageNumber,
		TotalPages:   totalPages,
		SearchQuery:  searchQuery,
	}

	err = tmpl.ExecuteTemplate(w, "listProducts.html", dataToSend)
	if err != nil {
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

	build := &data.Build{}

	getProduct := func(productType string, productID primitive.ObjectID) (data.Product, error) {
		products, _, err := productsListing(productType, bson.M{"_id": productID}, 0, 0, nil)
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

func getFullBuild(r *http.Request) (*data.FullBuild, error) {
	logger := logging.GetLogger()

	productHeaders, err := getBuildFromCookie(r)
	if err != nil {
		logger.Errorf("Failed to list build components: %v", err)
		return nil, err
	}

	var cpu data.Cpu
	var motherboard data.Motherboard
	var ram data.Ram
	var ssd data.Ssd
	var hdd data.Hdd
	var cooling data.Cooling
	var powersupply data.PowerSupply
	var gpu data.Gpu
	var housing data.Housing

	build := &data.FullBuild{
		CPU:         cpu,
		Motherboard: motherboard,
		RAM:         ram,
		SSD:         ssd,
		HDD:         hdd,
		Cooling:     cooling,
		PowerSupply: powersupply,
		Housing:     housing,
		GPU:         gpu,
	}

	for _, header := range productHeaders {

		switch header.ProductType {
		case "cpu":
			_, err = getAndDecodeProduct(data.CpuCollection, &cpu, header.ID)
			if err == nil {
				build.CPU = cpu
			}
		case "motherboard":
			_, err = getAndDecodeProduct(data.MotherboardCollection, &motherboard, header.ID)
			if err == nil {
				build.Motherboard = motherboard
			}
		case "ram":
			_, err = getAndDecodeProduct(data.RamCollection, &ram, header.ID)
			if err == nil {
				build.RAM = ram
			}
		case "gpu":
			_, err = getAndDecodeProduct(data.GpuCollection, &gpu, header.ID)
			if err == nil {
				build.GPU = gpu
			}
		case "ssd":
			_, err = getAndDecodeProduct(data.SsdCollection, &ssd, header.ID)
			if err == nil {
				build.SSD = ssd
			}
		case "hdd":
			_, err = getAndDecodeProduct(data.HddCollection, &hdd, header.ID)
			if err == nil {
				build.HDD = hdd
			}
		case "cooling":
			_, err = getAndDecodeProduct(data.CoolingCollection, &cooling, header.ID)
			if err == nil {
				build.Cooling = cooling
			}
		case "powersupply":
			_, err = getAndDecodeProduct(data.PowerSupplyCollection, &powersupply, header.ID)
			if err == nil {
				build.PowerSupply = powersupply
			}
		case "housing":
			_, err = getAndDecodeProduct(data.HousingCollection, &housing, header.ID)
			if err == nil {
				build.Housing = housing
			}
		default:
			logger.Errorf("Unknown product type: %s", header.ProductType)
			continue
		}

		if err != nil {
			logger.Errorf("Error getting product: %s, error: %v", header.ID, err)
			continue
		}
	}

	fmt.Println(build)
	return build, nil
}

func productsListing(productType string, filter bson.M, pageNumber int, perPageLimit int, w http.ResponseWriter) ([]data.Product, int, error) {
	logger := logging.GetLogger()
	skip := (pageNumber - 1) * perPageLimit

	collection, _, err := defineStruct(productType)

	if err != nil {
		return nil, 0, err
	}

	totalProducts, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		logger.Errorf("Error counting products: %v", err)
		return nil, 0, err
	}

	var productsList []data.Product

	cur, err := data.GetProducts(collection, filter, skip, perPageLimit)
	if err != nil {
		logger.Errorf("Error when trying to get products: %v", err)
		return nil, 0, err
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
				logger.Errorf("error: %v", err)
				continue
			}
			productsList = append(productsList, cpuItem.Standardize())
		case "cooling":
			var coolingItem data.Cooling
			err = cur.Decode(&coolingItem)
			if err != nil {
				logger.Errorf("error: %v", err)
				continue
			}
			productsList = append(productsList, coolingItem.Standardize())
		case "motherboard":
			var motherboardItem data.Motherboard
			err = cur.Decode(&motherboardItem)
			if err != nil {
				logger.Errorf("error: %v", err)
				continue
			}
			productsList = append(productsList, motherboardItem.Standardize())
		case "housing":
			var housingItem data.Housing
			err = cur.Decode(&housingItem)
			if err != nil {
				logger.Errorf("error: %v", err)
				continue
			}
			productsList = append(productsList, housingItem.Standardize())
		case "hdd":
			var hddItem data.Hdd
			err = cur.Decode(&hddItem)
			if err != nil {
				logger.Errorf("error: %v", err)
				continue
			}
			productsList = append(productsList, hddItem.Standardize())
		case "ssd":
			var ssdItem data.Ssd
			err = cur.Decode(&ssdItem)
			if err != nil {
				logger.Errorf("error: %v", err)
				continue
			}
			productsList = append(productsList, ssdItem.Standardize())
		case "powersupply":
			var powerSupplyItem data.PowerSupply
			err = cur.Decode(&powerSupplyItem)
			if err != nil {
				logger.Errorf("error: %v", err)
				continue
			}
			productsList = append(productsList, powerSupplyItem.Standardize())
		case "gpu":
			var gpuItem data.Gpu
			err = cur.Decode(&gpuItem)
			if err != nil {
				logger.Errorf("error: %v", err)
				continue
			}
			productsList = append(productsList, gpuItem.Standardize())
		case "ram":
			var ramItem data.Ram
			err = cur.Decode(&ramItem)
			if err != nil {
				logger.Errorf("error: %v", err)
				continue
			}
			productsList = append(productsList, ramItem.Standardize())
		default:
			logger.Errorf("Error: invalid product type value!")
			return productsList, int(totalProducts), errors.New("invalid product type")
		}
	}

	if err = cur.Err(); err != nil {
		logger.Error(err.Error())
	}

	logger.Infof("Found multiple items: %v", len(productsList))
	return productsList, int(totalProducts), nil
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

	collection, item, err := defineStruct(productType)
	if err != nil {
		HandleError(err, logger, w)
		return
	}

	ObjID, err := primitive.ObjectIDFromHex(r.FormValue("productID")[10:34])

	tmpl := template.Must(template.ParseFiles("html/productInformation.html"))

	item, err = getAndDecodeProduct(collection, item, ObjID)

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

func getAndDecodeProduct(collection *mongo.Collection, item interface{}, productID primitive.ObjectID) (interface{}, error) {
	filter := bson.M{"_id": productID}
	logger := logging.GetLogger()
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err = cur.Close(ctx)
		if err != nil {
			logger.Errorf("cursor error: %v", err)
		}
	}(cur, context.TODO())

	err = decodeProduct(cur, item)

	if err != nil {
		return nil, err
	}

	if err = cur.Err(); err != nil {
		logger.Errorf("error: %v", err)
	}

	logger.Infof("Found Item: %v", item)
	return item, nil
}

func getProductsFilter(productType string, r *http.Request) bson.M {
	var filter bson.M

	switch productType {
	case "cpu":
		filter = filters.FilterCpu(r)
	case "motherboard":
		filter = filters.FilterMotherboard(r)
	case "cooling":
		filter = filters.FilterCooling(r)
	case "housing":
		filter = filters.FilterHousing(r)
	case "powersupply":
		filter = filters.FilterPowerSupply(r)
	case "ssd":
		filter = filters.FilterSsd(r)
	case "hdd":
		filter = filters.FilterHdd(r)
	case "ram":
		filter = filters.FilterRam(r)
	case "gpu":
		filter = filters.FilterGpu(r)
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

func AddToBuild(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	productType := r.FormValue("productType")

	_, product, err := defineStruct(productType)
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

	collection, err := data.DefineCollection(productType)
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

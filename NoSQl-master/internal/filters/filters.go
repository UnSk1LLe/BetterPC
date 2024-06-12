package filters

import (
	"MongoDb/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

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

func getFloatInterval(from string, to string) (float64, float64, error) {
	logger := logging.GetLogger()
	var minValue, maxValue float64
	var err error
	if from != "" {
		minValue, err = strconv.ParseFloat(from, 64)
		if err != nil {
			logger.Errorf("Error parsing price from: %v", err)
			return minValue, maxValue, err
		}
	} else {
		minValue = 0
	}

	if to != "" {
		maxValue, err = strconv.ParseFloat(to, 64)
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

func getIntervalPrice(r *http.Request) (int, int, error) {
	keyFromStr := r.Form.Get("Price-min")
	keyToStr := r.Form.Get("Price-max")

	keyFrom := 0
	keyTo := 1000000
	var err error

	if keyFromStr != "" {
		keyFrom, err = strconv.Atoi(keyFromStr)
		if err != nil {
			return keyFrom, keyTo, err
		}
	}

	if keyToStr != "" {
		keyTo, err = strconv.Atoi(keyToStr)
		if err != nil {
			return keyFrom, keyTo, err
		}
	}

	return keyFrom, keyTo, nil
}

func getIntervalFilter(key string, r *http.Request) bson.M {
	keyFromStr := r.Form.Get(key + "-min")
	keyToStr := r.Form.Get(key + "-max")

	filter := bson.M{}

	keyFrom, KeyTo, err := getInterval(keyFromStr, keyToStr)
	if err != nil {
		return nil
	}

	if keyFrom != 0 || KeyTo != 0 {
		filter = bson.M{"$gte": keyFrom, "$lte": KeyTo}
	}
	return filter
}

func getFloatIntervalFilter(key string, r *http.Request) bson.M {
	keyFromStr := r.Form.Get(key + "-min")
	keyToStr := r.Form.Get(key + "-max")

	filter := bson.M{}

	keyFrom, KeyTo, err := getFloatInterval(keyFromStr, keyToStr)
	if err != nil {
		return nil
	}

	if keyFrom != 0 || KeyTo != 0 {
		filter = bson.M{"$gte": keyFrom, "$lte": KeyTo}
	}
	return filter
}

func getPriceFilter(r *http.Request) bson.M {
	logger := logging.GetLogger()
	filter := bson.M{}
	priceFrom, priceTo, err := getIntervalPrice(r)
	if err != nil {
		logger.Error("Could not get price interval")
		return filter
	}

	if priceFrom != 0 || priceTo != 0 {
		effectivePrice := bson.M{
			"$subtract": []interface{}{
				"$general.price",
				bson.M{
					"$multiply": []interface{}{
						"$general.price",
						bson.M{"$divide": []interface{}{"$general.discount", 100}},
					},
				},
			},
		}
		priceFilter := bson.M{}
		if priceFrom != 0 {
			priceFilter["$gte"] = priceFrom
		}
		if priceTo != 0 {
			priceFilter["$lte"] = priceTo
		}
		filter["$expr"] = bson.M{
			"$and": []bson.M{
				{"$gte": []interface{}{effectivePrice, priceFrom}},
				{"$lte": []interface{}{effectivePrice, priceTo}},
			},
		}
	}
	return filter
}

func SearchProducts(searchQuery string) bson.M {
	if searchQuery != "" {
		searchWords := strings.Split(searchQuery, " ")
		searchFilters := make([]bson.M, len(searchWords))
		for i, word := range searchWords {
			searchFilters[i] = bson.M{
				"$or": []bson.M{
					{"general.model": bson.M{"$regex": word, "$options": "i"}},
					{"general.manufacturer": bson.M{"$regex": word, "$options": "i"}},
					{"main.category": bson.M{"$regex": word, "$options": "i"}},
				},
			}
		}
		return bson.M{"$and": searchFilters}
	}
	return bson.M{}
}

func FilterCpu(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var categories []string
	var cores []int
	var ramTypes []string
	var sockets []string
	var year []int

	manufacturers = r.Form["Manufacturer"]
	categories = r.Form["Category"]
	coreValues := r.Form["P-Cores"]

	appendIntFilterParameters(coreValues, &cores)

	ramTypes = r.Form["Ram Type"]
	sockets = r.Form["Socket"]

	yearValues := r.Form["Year"]
	appendIntFilterParameters(yearValues, &year)

	filter := bson.M{}

	filter = getPriceFilter(r)

	filter["tdp"] = getIntervalFilter("TDP", r)

	integratedGraphicsValues := r.Form["Integrated graphics"]
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

	if len(year) > 0 {
		yearFilter := bson.M{"$in": year}
		filter["main.year"] = yearFilter
	}

	if len(categories) > 0 {
		categoryFilter := bson.M{"$in": categories}
		filter["main.category"] = categoryFilter
	}

	if len(cores) > 0 {
		coreFilter := bson.M{"$in": cores}
		filter["cores.p-cores"] = coreFilter
	}

	if len(ramTypes) > 0 {
		for _, t := range ramTypes {
			if t == "DDR4" {
				filter["ram.max_frequency.0"] = bson.M{"$gt": 0}
			}
			if t == "DDR5" {
				filter["ram.max_frequency.1"] = bson.M{"$gt": 0}
			}
		}
	}

	if len(sockets) > 0 {
		filter["main.socket"] = bson.M{"$in": sockets}
	}

	logger.Infof("Set filter: %v", filter)
	return filter
}

func FilterMotherboard(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var formFactors []string
	var chipsets []string
	var intelChipsets []string
	var amdChipsets []string
	var ramTypes []string
	var sockets []string
	var slotsValues []string
	var slots []int
	var m2InterfacesValues []string
	var m2Interfaces []int

	manufacturers = r.Form["Manufacturer"]
	formFactors = r.Form["Form-factor"]
	ramTypes = r.Form["RAM type"]
	sockets = r.Form["Socket"]
	intelChipsets = r.Form["Intel chipset groups"]
	amdChipsets = r.Form["AMD chipset groups"]
	slotsValues = r.Form["RAM slots"]
	m2InterfacesValues = r.Form["M2 interfaces"]

	appendIntFilterParameters(slotsValues, &slots)
	appendIntFilterParameters(m2InterfacesValues, &m2Interfaces)

	chipsets = append(chipsets, intelChipsets...)
	chipsets = append(chipsets, amdChipsets...)

	chipsetPattern := "^(" + strings.Join(chipsets, "|") + ")"

	chipsetFilter := bson.M{"$regex": primitive.Regex{Pattern: chipsetPattern, Options: "i"}}

	filter := bson.M{}
	filter = getPriceFilter(r)

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}
	if len(formFactors) > 0 {
		filter["form_factor"] = bson.M{"$in": formFactors}
	}
	if len(chipsets) > 0 {
		filter["chipset"] = chipsetFilter
	}
	if len(ramTypes) > 0 {
		filter["ram.type"] = bson.M{"$in": ramTypes}
	}
	if len(sockets) > 0 {
		filter["socket"] = bson.M{"$in": sockets}
	}
	if len(slots) > 0 {
		filter["ram.slots"] = bson.M{"$in": slots}
	}
	if len(m2Interfaces) > 0 {
		filter["interfaces.M2"] = bson.M{"$in": m2Interfaces}
	}

	logger.Infof("Set filter: %v", filter)
	return filter
}

func FilterPowerSupply(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var modularValues []string

	manufacturers = r.Form["Manufacturer"]
	modularValues = r.Form["Modular"]

	filter := bson.M{}

	filter = getPriceFilter(r)

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	modular := bson.M{}
	if len(modularValues) == 1 {
		value := modularValues[0]
		if value == "yes" {
			modular = bson.M{"$ne": false}
		} else if value == "no" {
			modular = bson.M{"$eq": false}
		}
		filter["modules"] = modular
	}

	filter["output_power"] = getIntervalFilter("Output power", r)

	logger.Infof("Set filter: %v", filter)
	return filter
}

func FilterCooling(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var types []string
	var sockets []string
	var mountTypes []string
	var fansNumberValues []string
	var fansNumber []int

	manufacturers = r.Form["Manufacturer"]
	types = r.Form["Type"]
	sockets = r.Form["Sockets"]
	mountTypes = r.Form["Mount Type"]
	fansNumberValues = r.Form["Fans"]
	appendIntFilterParameters(fansNumberValues, &fansNumber)

	filter := bson.M{}

	filter = getPriceFilter(r)

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(types) > 0 {
		filter["type"] = bson.M{"$in": types}
	}

	if len(sockets) > 0 {
		filter["sockets"] = bson.M{"$in": sockets}
	}

	if len(mountTypes) > 0 {
		filter["mount_type"] = bson.M{"$in": mountTypes}
	}

	if len(fansNumber) > 0 {
		filter["fans.0"] = bson.M{"$in": fansNumber}
	}

	filter["rpm.0"] = getIntervalFilter("RPM", r)
	filter["rpm.1"] = getIntervalFilter("RPM", r)
	filter["tdp"] = getIntervalFilter("TDP", r)
	filter["noise_level"] = getIntervalFilter("Noise Level", r)

	logger.Infof("Set filter: %v", filter)
	return filter
}

func FilterHousing(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var formFactors []string
	var driveBays35Values []string
	var driveBays35 []int
	var driveBays25Values []string
	var driveBays25 []int
	var mbFormFactors []string
	var expansionSlotsValues []string
	var expansionSlots []int

	manufacturers = r.Form["Manufacturer"]
	formFactors = r.Form["Form Factor"]
	driveBays35Values = r.Form["3.5 Drive Bays"]
	driveBays25Values = r.Form["2.5 Drive Bays"]
	mbFormFactors = r.Form["MB Form Factor"]
	expansionSlotsValues = r.Form["Expansion Slots"]

	appendIntFilterParameters(driveBays35Values, &driveBays35)
	appendIntFilterParameters(driveBays25Values, &driveBays25)
	appendIntFilterParameters(expansionSlotsValues, &expansionSlots)

	filter := bson.M{}

	filter = getPriceFilter(r)

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}
	if len(formFactors) > 0 {
		filter["form_factor"] = bson.M{"$in": formFactors}
	}
	if len(driveBays35) > 0 {
		filter["drive_bays.3_5"] = bson.M{"$in": driveBays35}
	}
	if len(driveBays25) > 0 {
		filter["drive_bays.2_5"] = bson.M{"$in": driveBays25}
	}
	if len(mbFormFactors) > 0 {
		filter["mb_form_factor"] = bson.M{"$in": mbFormFactors}
	}
	if len(expansionSlots) > 0 {
		filter["expansion_slots"] = bson.M{"$in": expansionSlots}
	}

	filter["graphic_card_size"] = getIntervalFilter("Graphic Card Size", r)
	filter["cooler_height"] = getIntervalFilter("Cooler Height", r)
	filter["weight"] = getIntervalFilter("Weight", r)

	logger.Infof("Set filter: %v", filter)
	return filter
}

func FilterRam(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var capacities []int
	var frequencies []int
	var types []string
	var formFactors []string
	var voltagesValues []string
	var voltages []float64
	var casLatencies []string

	manufacturers = r.Form["Manufacturer"]
	capacityValues := r.Form["Capacity"]
	appendIntFilterParameters(capacityValues, &capacities)

	frequencyValues := r.Form["Frequency"]
	appendIntFilterParameters(frequencyValues, &frequencies)

	types = r.Form["Type"]
	formFactors = r.Form["Form-factor"]

	voltagesValues = r.Form["Voltage"]
	appendFloatFilterParameters(voltagesValues, &voltages)

	casLatencies = r.Form["CAS Latency"]

	filter := bson.M{}

	filter = getPriceFilter(r)

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(capacities) > 0 {
		filter["capacity"] = bson.M{"$in": capacities}
	}

	filter["frequency"] = getIntervalFilter("Frequency", r)

	if len(types) > 0 {
		filter["type"] = bson.M{"$in": types}
	}

	if len(formFactors) > 0 {
		filter["form_factor"] = bson.M{"$in": formFactors}
	}

	if len(casLatencies) > 0 {
		filter["cas_latency"] = bson.M{"$in": casLatencies}
	}

	filter["voltage"] = getFloatIntervalFilter("Voltage", r)

	logger.Infof("Set RAM filter: %v", filter)
	return filter
}

func FilterHdd(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var capacities []int
	var interfaces []string
	var writeMethods []string
	var spindleSpeeds []int
	var formFactors []string

	manufacturers = r.Form["Manufacturer"]
	capacityValues := r.Form["Capacity"]
	appendIntFilterParameters(capacityValues, &capacities)

	interfaces = r.Form["Interface"]
	writeMethods = r.Form["Write Method"]
	spindleSpeedValues := r.Form["Spindle Speed"]
	appendIntFilterParameters(spindleSpeedValues, &spindleSpeeds)

	formFactors = r.Form["FormFactor"]

	filter := bson.M{}

	filter = getPriceFilter(r)

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(capacities) > 0 {
		filter["capacity"] = bson.M{"$in": capacities}
	}

	if len(interfaces) > 0 {
		filter["interface"] = bson.M{"$in": interfaces}
	}

	if len(writeMethods) > 0 {
		filter["write_method"] = bson.M{"$in": writeMethods}
	}

	if len(spindleSpeeds) > 0 {
		filter["spindle_speed"] = bson.M{"$in": spindleSpeeds}
	}

	filter["transfer_rate"] = getIntervalFilter("Transfer Rate", r)

	if len(formFactors) > 0 {
		filter["form_factor"] = bson.M{"$in": formFactors}
	}

	logger.Infof("Set HDD filter: %v", filter)
	return filter
}

func FilterSsd(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var capacities []int
	var memoryTypes []string
	var formFactors []string

	manufacturers = r.Form["Manufacturer"]
	capacityValues := r.Form["Capacity"]
	appendIntFilterParameters(capacityValues, &capacities)

	memoryTypes = r.Form["Memory Type"]

	formFactors = r.Form["Form Factor"]

	filter := bson.M{}

	filter = getPriceFilter(r)

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(capacities) > 0 {
		filter["capacity"] = bson.M{"$in": capacities}
	}

	if len(memoryTypes) > 0 {
		regexPattern := strings.Join(memoryTypes, "|")
		filter["memory_type"] = bson.M{"$regex": regexPattern, "$options": "i"}
	}

	filter["read"] = getIntervalFilter("Read Speed", r)
	filter["write"] = getIntervalFilter("Write Speed", r)

	if len(formFactors) > 0 {
		escapedFormFactors := make([]string, len(formFactors))
		for i, formFactor := range formFactors {
			escapedFormFactors[i] = "\\b" + regexp.QuoteMeta(formFactor) + "\\b"
		}
		regexPattern := strings.Join(escapedFormFactors, "|")
		filter["form_factor"] = bson.M{"$regex": regexPattern, "$options": "i"}
	}

	logger.Infof("Set SSD filter: %v", filter)
	return filter
}

func FilterGpu(r *http.Request) bson.M {
	logger := logging.GetLogger()

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Could not parse form: %v", err)
		return nil
	}

	var manufacturers []string
	var architectures []string
	var memoryCapacityValues []string
	var memoryCapacity []int
	var memoryTypes []string
	var technicalProcessValues []string
	var technicalProcess []int
	var maxResolution []string
	var maxMonitors []int
	var maxMonitorsValues []string

	manufacturers = r.Form["Manufacturer"]
	architectures = r.Form["Architecture"]
	memoryCapacityValues = r.Form["Memory capacity"]
	memoryTypes = r.Form["Memory type"]
	maxResolution = r.Form["Max resolution"]
	maxMonitorsValues = r.Form["Max monitors"]

	appendIntFilterParameters(memoryCapacityValues, &memoryCapacity)
	appendIntFilterParameters(maxMonitorsValues, &maxMonitors)
	appendIntFilterParameters(technicalProcessValues, &technicalProcess)

	filter := bson.M{}

	filter = getPriceFilter(r)

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(architectures) > 0 {
		filter["architecture"] = bson.M{"$in": architectures}
	}

	if len(memoryCapacity) > 0 {
		filter["memory.capacity"] = bson.M{"$in": memoryCapacity}
	}

	if len(memoryTypes) > 0 {
		filter["memory.type"] = bson.M{"$in": memoryTypes}
	}

	filter["gpu_frequency"] = getIntervalFilter("GPU frequency", r)

	if len(technicalProcess) > 0 {
		filter["process_size"] = bson.M{"$in": technicalProcess}
	}

	filter["tdp"] = getIntervalFilter("TDP", r)

	if len(maxResolution) > 0 {
		filter["max_resolution"] = bson.M{"$in": maxResolution}
	}

	if len(maxMonitors) > 0 {
		filter["max_monitors"] = bson.M{"$in": maxMonitors}
	}

	logger.Infof("Set GPU filter: %v", filter)
	return filter
}

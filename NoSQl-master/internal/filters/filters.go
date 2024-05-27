package filters

import (
	"MongoDb/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
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
	var ramTypes []string
	var sockets []string
	var pcie []string

	manufacturers = r.Form["Manufacturer"]
	formFactors = r.Form["Form-factor"]
	chipsets = r.Form["Chipset"]
	ramTypes = r.Form["Ram-type"]
	sockets = r.Form["Socket"]
	pcie = r.Form["PCI-E"]

	priceFromStr := r.Form.Get("Price-min")
	priceToStr := r.Form.Get("Price-max")

	priceFrom, priceTo, err := getInterval(priceFromStr, priceToStr)

	filter := bson.M{}

	if priceFrom != 0 || priceTo != 0 {
		filter["general.price"] = bson.M{"$gte": priceFrom, "$lte": priceTo}
	}

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(formFactors) > 0 {
		filter["form_factor"] = bson.M{"$in": formFactors}
	}

	if len(chipsets) > 0 {
		filter["chipset"] = bson.M{"$in": chipsets}
	}

	if len(ramTypes) > 0 {
		filter["ram.type"] = bson.M{"$in": ramTypes}
	}

	if len(sockets) > 0 {
		filter["socket"] = bson.M{"$in": sockets}
	}

	if len(pcie) > 0 {
		filter["pci_standard"] = bson.M{"$in": pcie}
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
	var formFactors []string
	var modular []string
	var gpuPower []string
	var cpuPower []string

	manufacturers = r.Form["Manufacturer"]
	formFactors = r.Form["Form-factor"]
	modular = r.Form["Modular"]
	gpuPower = r.Form["GPU-power"]
	cpuPower = r.Form["CPU-power"]

	priceFromStr := r.Form.Get("Price-min")
	priceToStr := r.Form.Get("Price-max")

	priceFrom, priceTo, err := getInterval(priceFromStr, priceToStr)

	filter := bson.M{}

	if priceFrom != 0 || priceTo != 0 {
		filter["general.price"] = bson.M{"$gte": priceFrom, "$lte": priceTo}
	}

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(formFactors) > 0 {
		filter["form_factor"] = bson.M{"$in": formFactors}
	}

	if len(modular) > 0 {
		filter["modules"] = bson.M{"$in": modular}
	}

	if len(gpuPower) > 0 {
		filter["connectors.PCI_E"] = bson.M{"$in": gpuPower}
	}

	if len(cpuPower) > 0 {
		filter["cpu_power"] = bson.M{"$in": cpuPower}
	}

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

	manufacturers = r.Form["Manufacturer"]
	types = r.Form["Type"]
	sockets = r.Form["Sockets"]
	mountTypes = r.Form["Mount-Type"]

	priceFromStr := r.Form.Get("Price-min")
	priceToStr := r.Form.Get("Price-max")

	priceFrom, priceTo, err := getInterval(priceFromStr, priceToStr)

	filter := bson.M{}

	if priceFrom != 0 || priceTo != 0 {
		filter["general.price"] = bson.M{"$gte": priceFrom, "$lte": priceTo}
	}

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
	var mbFormFactors []string
	var psFormFactors []string

	manufacturers = r.Form["Manufacturer"]
	formFactors = r.Form["Form-Factor"]
	mbFormFactors = r.Form["MB-Form-Factor"]
	psFormFactors = r.Form["PS-Form-Factor"]

	priceFromStr := r.Form.Get("Price-min")
	priceToStr := r.Form.Get("Price-max")

	priceFrom, priceTo, err := getInterval(priceFromStr, priceToStr)

	filter := bson.M{}

	if priceFrom != 0 || priceTo != 0 {
		filter["general.price"] = bson.M{"$gte": priceFrom, "$lte": priceTo}
	}

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(formFactors) > 0 {
		filter["form_factor"] = bson.M{"$in": formFactors}
	}

	if len(mbFormFactors) > 0 {
		filter["mb_form_factor"] = bson.M{"$in": mbFormFactors}
	}

	if len(psFormFactors) > 0 {
		filter["ps_form_factor"] = bson.M{"$in": psFormFactors}
	}

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

	manufacturers = r.Form["Manufacturer"]
	capacityValues := r.Form["Capacity"]
	appendIntFilterParameters(capacityValues, &capacities)

	frequencyValues := r.Form["Frequency"]
	appendIntFilterParameters(frequencyValues, &frequencies)

	types = r.Form["ram-type"]
	formFactors = r.Form["form-factor"]

	filter := bson.M{}

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(capacities) > 0 {
		filter["capacity"] = bson.M{"$in": capacities}
	}

	if len(frequencies) > 0 {
		filter["frequency"] = bson.M{"$in": frequencies}
	}

	if len(types) > 0 {
		filter["type"] = bson.M{"$in": types}
	}

	if len(formFactors) > 0 {
		filter["form_factor"] = bson.M{"$in": formFactors}
	}

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
	writeMethods = r.Form["WriteMethod"]
	spindleSpeedValues := r.Form["SpindleSpeed"]
	appendIntFilterParameters(spindleSpeedValues, &spindleSpeeds)

	formFactors = r.Form["FormFactor"]

	filter := bson.M{}

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
	var interfaces []string
	var memoryTypes []string
	var formFactors []string

	manufacturers = r.Form["Manufacturer"]
	capacityValues := r.Form["Capacity"]
	appendIntFilterParameters(capacityValues, &capacities)

	interfaces = r.Form["Interface"]
	memoryTypes = r.Form["MemoryType"]

	formFactors = r.Form["FormFactor"]

	filter := bson.M{}

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(capacities) > 0 {
		filter["capacity"] = bson.M{"$in": capacities}
	}

	if len(interfaces) > 0 {
		filter["interface"] = bson.M{"$in": interfaces}
	}

	if len(memoryTypes) > 0 {
		filter["memory_type"] = bson.M{"$in": memoryTypes}
	}

	if len(formFactors) > 0 {
		filter["form_factor"] = bson.M{"$in": formFactors}
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
	var memoryTypes []string
	var interfaces []string
	var coolingTypes []string

	manufacturers = r.Form["Manufacturer"]
	architectures = r.Form["Architecture"]
	memoryTypes = r.Form["MemoryType"]
	interfaces = r.Form["Interface"]
	coolingTypes = r.Form["CoolingType"]

	filter := bson.M{}

	if len(manufacturers) > 0 {
		filter["general.manufacturer"] = bson.M{"$in": manufacturers}
	}

	if len(architectures) > 0 {
		filter["architecture"] = bson.M{"$in": architectures}
	}

	if len(memoryTypes) > 0 {
		filter["memory.type"] = bson.M{"$in": memoryTypes}
	}

	if len(interfaces) > 0 {
		filter["interfaces.type"] = bson.M{"$in": interfaces}
	}

	if len(coolingTypes) > 0 {
		filter["cooling.type"] = bson.M{"$in": coolingTypes}
	}

	logger.Infof("Set GPU filter: %v", filter)
	return filter
}

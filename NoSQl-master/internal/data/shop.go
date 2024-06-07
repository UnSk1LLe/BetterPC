package data

import (
	"MongoDb/pkg/logging"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"reflect"
	"strconv"
	"strings"
)

var CpuCollection *mongo.Collection
var MotherboardCollection *mongo.Collection
var RamCollection *mongo.Collection
var GpuCollection *mongo.Collection
var CoolingCollection *mongo.Collection
var SsdCollection *mongo.Collection
var HddCollection *mongo.Collection
var HousingCollection *mongo.Collection
var PowerSupplyCollection *mongo.Collection

type ProductHeader struct {
	ID          primitive.ObjectID
	ProductType string
}

type Product struct {
	ProductHeader ProductHeader
	General       General
	Name          string
	Description   string
}

type Build struct {
	CPU         Product
	Motherboard Product
	RAM         Product
	GPU         Product
	SSD         Product
	HDD         Product
	Cooling     Product
	PowerSupply Product
	Housing     Product
}

type FullBuild struct {
	CPU         Cpu
	Motherboard Motherboard
	RAM         Ram
	GPU         Gpu
	SSD         Ssd
	HDD         Hdd
	Cooling     Cooling
	PowerSupply PowerSupply
	Housing     Housing
}

type Cpu struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	General        General            `bson:"general"`
	Main           MainCpu            `bson:"main"`
	Cores          CoresCpu           `bson:"cores"`
	ClockFrequency ClockFrequencyCpu  `bson:"clock_frequency" json:"clock_frequency,omitempty"`
	Ram            RamCpu             `bson:"ram"`
	Tdp            int                `bson:"tdp" json:"tdp,omitempty"`
	Graphics       string             `bson:"graphics" json:"graphics,omitempty"`
	PciE           int                `bson:"pci-e" json:"pci_e,omitempty"`
	MaxTemperature int                `bson:"max_temperature"`
}

type MainCpu struct {
	Category   string `bson:"category" json:"category,omitempty"`
	Generation string `bson:"generation" json:"generation,omitempty"`
	Socket     string `bson:"socket" json:"socket,omitempty"`
	Year       int    `bson:"year" json:"year,omitempty"`
}

type CoresCpu struct {
	Pcores           int `bson:"p-cores" json:"p-cores,omitempty"`
	Ecores           int `bson:"e-cores" json:"e-cores,omitempty"`
	Threads          int `bson:"threads" json:"threads,omitempty"`
	TechnicalProcess int `bson:"technical_process" json:"technical_process,omitempty"`
}

type ClockFrequencyCpu struct {
	Pcores         []float64 `bson:"p-cores" json:"p-cores,omitempty"`
	Ecores         []float64 `bson:"e-cores" json:"e-cores,omitempty"`
	FreeMultiplier bool      `bson:"free_multiplier" json:"free_multiplier,omitempty"`
}

type RamCpu struct {
	Channels     int   `bson:"channels"`
	MaxFrequency []int `bson:"max_frequency"`
	MaxCapacity  int   `bson:"max_capacity"`
}

type Motherboard struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	General     General            `bson:"general"`
	Socket      string             `bson:"socket"`
	Chipset     string             `bson:"chipset"`
	FormFactor  string             `bson:"form_factor"`
	Ram         ramMb              `bson:"ram"`
	Interfaces  interfaces         `bson:"interfaces"`
	PciStandard int                `bson:"pci_standard"`
	MbPower     int                `bson:"mb_power"`
	CpuPower    int                `bson:"cpu_power"`
}

type ramMb struct {
	Slots        int    `bson:"slots"`
	Type         string `bson:"type"`
	MaxFrequency int    `bson:"max_frequency"`
	MaxCapacity  int    `bson:"max_capacity"`
}

type interfaces struct {
	Sata3   int `bson:"SATA3"`
	M2      int `bson:"M2"`
	PciE1x  int `bson:"PCI_E_1x"`
	PciE16x int `bson:"PCI_E_16x"`
}

type Ram struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	General      General            `bson:"general"`
	Capacity     int                `bson:"capacity"`
	Number       int                `bson:"number"`
	FormFactor   string             `bson:"form_factor"`
	Rank         int                `bson:"rank"`
	Type         string             `bson:"type"`
	Frequency    int                `bson:"frequency"`
	Bandwidth    int                `bson:"bandwidth"`
	CasLatency   string             `bson:"cas_latency"`
	TimingScheme []int              `bson:"timing_scheme"`
	Voltage      float64            `bson:"voltage"`
	Cooling      string             `bson:"cooling"`
	Height       int                `bson:"height"`
}

type Ssd struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	General    General            `bson:"general"`
	Type       string             `bson:"type"`
	Capacity   int                `bson:"capacity"`
	Interface  string             `bson:"interface"`
	MemoryType string             `bson:"memory_type"`
	Read       int                `bson:"read"`
	Write      int                `bson:"write"`
	FormFactor string             `bson:"form_factor"`
	Mftb       float64            `bson:"mftb"`
	Size       []float64          `bson:"size"`
	Weight     int                `bson:"weight"`
}

type Hdd struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	General      General            `bson:"general"`
	Type         string             `bson:"type"`
	Capacity     int                `bson:"capacity"`
	Interface    string             `bson:"interface"`
	WriteMethod  string             `bson:"write_method"`
	TransferRate int                `bson:"transfer_rate"`
	SpindleSpeed int                `bson:"spindle_speed"`
	FormFactor   string             `bson:"form_factor"`
	Mftb         int                `bson:"mftb"`
	Size         []float64          `bson:"size"`
	Weight       int                `bson:"weight"`
}

type Gpu struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	General       General            `bson:"general"`
	Architecture  string             `bson:"architecture"`
	Memory        memoryGpu          `bson:"memory"`
	GpuFrequency  int                `bson:"gpu_frequency"`
	ProcessSize   int                `bson:"process_size"`
	MaxResolution string             `bson:"max_resolution"`
	Interfaces    []interfacesGpu    `bson:"interfaces"`
	MaxMonitors   int                `bson:"max_monitors"`
	Cooling       coolingGpu         `bson:"cooling"`
	Tdp           int                `bson:"tdp"`
	TdpR          int                `bson:"tdp_r"`
	PowerSupply   []int              `bson:"power_supply"`
	Slots         float64            `bson:"slots"`
	Size          []int              `bson:"size"`
}

type memoryGpu struct {
	Capacity       int    `bson:"capacity"`
	Type           string `bson:"type"`
	InterfaceWidth int    `bson:"interface_width"`
	Frequency      int    `bson:"frequency"`
}

type interfacesGpu struct {
	Type   string `bson:"type"`
	Number int    `bson:"number"`
}

type coolingGpu struct {
	Type      string `bson:"type"`
	FanNumber int    `bson:"fan_number"`
}

type Cooling struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	General    General            `bson:"general"`
	Type       string             `bson:"type"`
	Sockets    []string           `bson:"sockets"`
	Fans       []int              `bson:"fans"`
	Rpm        []int              `bson:"rpm"`
	Tdp        int                `bson:"tdp"`
	NoiseLevel int                `bson:"noise_level"`
	MountType  string             `bson:"mount_type"`
	Power      int                `bson:"power"`
	Height     int                `bson:"height"`
}

type Housing struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	General         General            `bson:"general"`
	FormFactor      string             `bson:"form_factor"`
	DriveBays       driveBays          `bson:"drive_bays"`
	MbFormFactor    string             `bson:"mb_form_factor"`
	PsFormFactor    string             `bson:"ps_form_factor"`
	ExpansionSlots  int                `bson:"expansion_slots"`
	GraphicCardSize int                `bson:"graphic_card_size"`
	CoolerHeight    int                `bson:"cooler_height"`
	Size            []int              `bson:"size"`
	Weight          float64            `bson:"weight"`
}

type driveBays struct {
	D35 int `bson:"3_5"`
	D25 int `bson:"2_5"`
}

type PowerSupply struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	General     General            `bson:"general"`
	FormFactor  string             `bson:"form_factor"`
	OutputPower int                `bson:"output_power"`
	Connectors  connectors         `bson:"connectors"`
	Modules     bool               `bson:"modules"`
	MbPower     int                `bson:"mb_power"`
	CpuPower    bsoncore.Array     `bson:"cpu_power"`
}

type connectors struct {
	Sata  int   `bson:"SATA"`
	Molex int   `bson:"MOLEX"`
	PciE  []int `bson:"PCI_E"`
}

type General struct {
	Manufacturer string `bson:"manufacturer"`
	Model        string `bson:"model"`
	Price        int    `bson:"price"`
	Discount     int    `bson:"discount"`
	Amount       int    `bson:"amount"`
}

func (g General) ProductFinalPrice() int {
	return g.Price - (g.Price * g.Discount / 100)
}

func GetProductById(collection *mongo.Collection, ID primitive.ObjectID) (*mongo.SingleResult, error) {
	logger := logging.GetLogger()
	result := collection.FindOne(context.TODO(), bson.M{"_id": ID})
	if result.Err() != nil {
		logger.Errorf(result.Err().Error())
		return nil, result.Err()
	}
	logger.Infof("found product with ID <%v> in <%v> collection", ID, collection.Name())
	return result, nil
}

func DeleteProductById(productType string, ID primitive.ObjectID) (*mongo.DeleteResult, error) {
	logger := logging.GetLogger()
	collection, err := DefineCollection(productType)
	if err != nil {
		logger.Errorf("Error deleting product, could not define type <%s>): %v", productType, err)
		return nil, err
	}
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": ID})
	if err != nil {
		logger.Errorf("Could not delete product of type(%s): %v", productType, err)
		return nil, err
	}
	logger.Infof("Product <%s> with ID: %v was DELETED!", productType, ID)
	return result, err
}

func (cpu Cpu) Standardize() Product {
	var product Product
	product.ProductHeader.ID = cpu.ID
	product.ProductHeader.ProductType = "cpu"
	product.General = cpu.General
	product.Name = cpu.Main.Category + " " + cpu.General.Model
	var cores string
	if cpu.Cores.Ecores > 0 {
		cores = "P-cores: " + strconv.Itoa(cpu.Cores.Pcores) + " E-cores: " + strconv.Itoa(cpu.Cores.Ecores) + ","
	} else {
		cores = strconv.Itoa(cpu.Cores.Pcores) + ","
	}
	product.Description = cpu.Main.Category + ", " + cpu.Main.Generation + " Generation, " +
		cpu.Main.Socket + " Socket, " + "Cores: " + cores + " Threads: " + strconv.Itoa(cpu.Cores.Threads) +
		", Technical process " + strconv.Itoa(cpu.Cores.TechnicalProcess) + " nm, "
	return product
}

func (motherboard Motherboard) Standardize() Product {
	var product Product
	product.ProductHeader.ID = motherboard.ID
	product.ProductHeader.ProductType = "motherboard"
	product.Name = motherboard.General.Model
	product.General = motherboard.General
	product.Description = "Socket: " + motherboard.Socket + ", Chipset: " + motherboard.Chipset +
		", Form Factor: " + motherboard.FormFactor + ", RAM: " + motherboard.Ram.Type + " " +
		strconv.Itoa(motherboard.Ram.MaxCapacity) + "GB"
	return product
}

func (cooling Cooling) Standardize() Product {
	var product Product
	product.ProductHeader.ID = cooling.ID
	product.ProductHeader.ProductType = "cooling"
	product.Name = cooling.General.Model
	product.General = cooling.General
	product.Description = "Type: " + cooling.Type + ", Sockets: " + strings.Join(cooling.Sockets, ", ") +
		", Fans: " + strconv.Itoa(len(cooling.Fans)) + ", TDP: " + strconv.Itoa(cooling.Tdp) + "W"
	return product
}

func (ram Ram) Standardize() Product {
	var product Product
	product.ProductHeader.ID = ram.ID
	product.ProductHeader.ProductType = "ram"
	product.Name = ram.General.Model
	product.General = ram.General
	product.Description = "Capacity: " + strconv.Itoa(ram.Capacity) + "GB, Type: " + ram.Type +
		", Frequency: " + strconv.Itoa(ram.Frequency) + "MHz, CAS Latency: " + ram.CasLatency
	return product
}

func (ssd Ssd) Standardize() Product {
	var product Product
	product.ProductHeader.ID = ssd.ID
	product.ProductHeader.ProductType = "ssd"
	product.Name = ssd.General.Model
	product.General = ssd.General
	product.Description = "Type: " + ssd.Type + ", Capacity: " + strconv.Itoa(ssd.Capacity) + "GB, " +
		"Interface: " + ssd.Interface + ", Read Speed: " + strconv.Itoa(ssd.Read) + "MB/s, " +
		"Write Speed: " + strconv.Itoa(ssd.Write) + "MB/s"
	return product
}

func (hdd Hdd) Standardize() Product {
	var product Product
	product.ProductHeader.ID = hdd.ID
	product.ProductHeader.ProductType = "hdd"
	product.Name = hdd.General.Model
	product.General = hdd.General
	product.Description = "Type: " + hdd.Type + ", Capacity: " + strconv.Itoa(hdd.Capacity) + "GB, " +
		"Interface: " + hdd.Interface + ", Spindle Speed: " + strconv.Itoa(hdd.SpindleSpeed) + "RPM"
	return product
}

func (gpu Gpu) Standardize() Product {
	var product Product
	product.ProductHeader.ID = gpu.ID
	product.ProductHeader.ProductType = "gpu"
	product.Name = gpu.General.Model
	product.General = gpu.General
	product.Description = "Architecture: " + gpu.Architecture + ", Memory: " + strconv.Itoa(gpu.Memory.Capacity) +
		"GB " + gpu.Memory.Type + ", Frequency: " + strconv.Itoa(gpu.GpuFrequency) + "MHz, " +
		"Max Resolution: " + gpu.MaxResolution
	return product
}

func (powerSupply PowerSupply) Standardize() Product {
	var product Product
	product.ProductHeader.ID = powerSupply.ID
	product.ProductHeader.ProductType = "powersupply"
	product.Name = powerSupply.General.Model
	product.General = powerSupply.General
	product.Description = "Form Factor: " + powerSupply.FormFactor + ", Output Power: " +
		strconv.Itoa(powerSupply.OutputPower) + "W, Modular: " + strconv.FormatBool(powerSupply.Modules)
	return product
}

func (housing Housing) Standardize() Product {
	var product Product
	product.ProductHeader.ID = housing.ID
	product.ProductHeader.ProductType = "housing"
	product.Name = housing.General.Model
	product.General = housing.General
	product.Description = "Form Factor: " + housing.FormFactor + ", Motherboard Form Factor: " +
		housing.MbFormFactor + ", Expansion Slots: " + strconv.Itoa(housing.ExpansionSlots)
	return product
}

func IsZero(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

func GetProducts(collection *mongo.Collection, filter bson.M, skip int, limit int) (*mongo.Cursor, error) {
	//if skip or limit not needed then set to 0
	logger := logging.GetLogger()
	findOptions := options.Find()
	if skip > 0 {
		findOptions.SetSkip(int64(skip))
	}
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		logger.Errorf("Error getting products from <%s> collection: %v", collection.Name(), err)
		return nil, err
	}
	logger.Infof("Found products from <%s> collection", collection.Name())
	return cur, nil
}

func UpdateProduct(collection *mongo.Collection, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	logger := logging.GetLogger()
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Errorf("Error updating product from <%s> collection: %v", collection.Name(), err)
		return nil, err
	}
	logger.Infof("Product from <%s> collection was UPDATED!", collection.Name())
	return result, nil
}

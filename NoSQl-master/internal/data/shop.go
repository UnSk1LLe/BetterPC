package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type CompatibilityMode struct {
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
	Channels     int    `bson:"channels"`
	Type         string `bson:"type"`
	MaxFrequency int    `bson:"max_frequency"`
	MaxCapacity  int    `bson:"max_capacity"`
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

func GetProductById(dbName string, collectionName string, ID primitive.ObjectID) (*mongo.SingleResult, error) {
	err := Init(dbName, collectionName)
	if err != nil {
		return nil, err
	}
	defer CloseConnection()

	result := Collection.FindOne(context.TODO(), bson.M{"_id": ID})
	return result, nil
}

func GetProducts(dbName string, collectionName string, filter bson.M) (*mongo.Cursor, error) {
	err := Init(dbName, collectionName)
	if err != nil {
		return nil, err
	}
	defer CloseConnection()
	cur, err := Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return cur, nil
}

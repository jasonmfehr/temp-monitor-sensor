package sensoriface

type Data interface {
	GetData() ([]DataPoint, error)
}

type DataPoint struct {
	SensorID    string
	Temperature float32
	Unit        TemperatureUnit
	Error       error
}

type TemperatureUnit string

const (
	Fahrenheit TemperatureUnit = "F"
	Celsius    TemperatureUnit = "C"
)

package raspberry_pi

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/jasonmfehr/temp-monitor-sensor/pkg/sensoriface"
)

// ioutilReadFile enables mocking ioutil.ReadFile
var ioutilReadFile = ioutil.ReadFile

// funcReadSensor enables mocking readSensor
var funcReadSensor = readSensor

// regexpCompile enables mocking regexp.Compile
var regexpCompile = regexp.Compile

func (s *DallasOneWire) GetData() []sensoriface.DataPoint {
	sensorIDs, err := ioutilReadFile(sensorIDsFile)
	if err != nil {
		return []sensoriface.DataPoint{
			sensoriface.DataPoint{
				Error: errors.Wrapf(err, "could not read file '%s'", sensorIDsFile),
			},
		}
	}

	dataPoints := []sensoriface.DataPoint{}
	for _, sensorID := range strings.Split(string(sensorIDs), "\n") {
		sensorID := strings.TrimSpace(sensorID)
		// TODO - is this if statement necessary?
		if len(sensorID) > 0 {
			dataPoints = append(dataPoints, funcReadSensor(sensorID))
		}
	}

	return dataPoints
}

func readSensor(sensorID string) sensoriface.DataPoint {
	dataFile := fmt.Sprintf("%s/%s/%s", sensorRootPath, sensorID, sensorDataFile)
	sensorData, err := ioutilReadFile(dataFile)
	if err != nil {
		return sensoriface.DataPoint{
			SensorID: sensorID,
			Error:    errors.Wrapf(err, "could not read sensor data from file '%s'", dataFile),
		}
	}

	sensorDataLines := strings.Split(string(sensorData), "\n")
	if len(sensorDataLines) < 2 {
		return sensoriface.DataPoint{
			SensorID: sensorID,
			Error:    errors.Errorf("sensor data has unexpected number of lines '%d'", len(sensorDataLines)),
		}
	}

	re, err := regexpCompile(tempPattern)
	if err != nil {
		return sensoriface.DataPoint{
			SensorID: sensorID,
			Error:    errors.Wrapf(err, "could not compile temperature pattern regular expression '%s'", tempPattern),
		}
	}

	matches := re.FindStringSubmatch(sensorDataLines[1])
	if len(matches) != 2 {
		return sensoriface.DataPoint{
			SensorID: sensorID,
			Error:    errors.Errorf("temp data did not match expected pattern '%s'", tempPattern),
		}
	}

	tempInt, err := strconv.Atoi(matches[1])
	if err != nil {
		return sensoriface.DataPoint{
			SensorID: sensorID,
			Error:    errors.Wrapf(err, "could not convert temperature reading '%s' into integer", matches[1]),
		}
	}

	return sensoriface.DataPoint{
		SensorID:    sensorID,
		Temperature: float32(tempInt) / 1000,
		Unit:        sensoriface.Celsius,
	}
}

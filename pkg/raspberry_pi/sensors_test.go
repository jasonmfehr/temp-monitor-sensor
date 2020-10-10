package raspberry_pi

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/jasonmfehr/temp-monitor-sensor/pkg/sensoriface"
	"github.com/stretchr/testify/assert"
)

func TestGetDataHappyPath(t *testing.T) {
	funcReadSensorCallCount := 0
	testSensor0 := sensoriface.DataPoint{
		SensorID: "sensor-0",
	}
	testSensor1 := sensoriface.DataPoint{
		SensorID: "sensor-1",
	}

	ioutilReadFile = func(filename string) ([]byte, error) {
		assert.Equal(t, "/sys/bus/w1/devices/w1_bus_master1/w1_master_slaves", filename)
		return []byte(fmt.Sprintf(" %s \n%s\n", testSensor0.SensorID, testSensor1.SensorID)), nil
	}

	funcReadSensor = func(sensorID string) sensoriface.DataPoint {
		funcReadSensorCallCount++
		if sensorID == testSensor0.SensorID {
			return testSensor0
		} else if sensorID == testSensor1.SensorID {
			return testSensor1
		}

		t.Errorf("unexpected sensorID: '%s'", sensorID)
		return sensoriface.DataPoint{}
	}

	actual := (&DallasOneWire{}).GetData()

	assert.Len(t, actual, 2)
	assert.Contains(t, actual, testSensor0)
	assert.Contains(t, actual, testSensor1)
	assert.Equal(t, 2, funcReadSensorCallCount)
}

func TestGetDataReadError(t *testing.T) {
	testErr := "test-get-data-read-error"

	ioutilReadFile = func(filename string) ([]byte, error) {
		return nil, fmt.Errorf(testErr)
	}

	actual := (&DallasOneWire{}).GetData()

	assert.Len(t, actual, 1)
	assert.EqualError(t, actual[0].Error, fmt.Sprintf("could not read file '%s': %s", sensorIDsFile, testErr))
}

func TestReadSensorHappyPath(t *testing.T) {
	testSensorID := "test-sensor"
	ioutilReadFile = func(filename string) ([]byte, error) {
		return []byte("5c 01 4b 46 7f ff 0c 10 d7 : crc=d7 YES\n5c 01 4b 46 7f ff 0c 10 d7 t=21750\n"), nil
	}

	regexpCompile = regexp.Compile

	actual := readSensor(testSensorID)

	assert.Equal(t, testSensorID, actual.SensorID)
	assert.Equal(t, float32(21.75), actual.Temperature)
	assert.Equal(t, sensoriface.Celsius, actual.Unit)
	assert.NoError(t, actual.Error)
}

func TestReadSensorFileError(t *testing.T) {
	testErr := "test-read-sensor-file-error"
	testSensorID := "foo-sensor"

	ioutilReadFile = func(filename string) ([]byte, error) {
		return nil, fmt.Errorf(testErr)
	}

	actual := readSensor(testSensorID)

	assert.EqualError(t, actual.Error, fmt.Sprintf("could not read sensor data from file '/sys/bus/w1/devices/%s/w1_slave': %s", testSensorID, testErr))
	assert.Equal(t, testSensorID, actual.SensorID)
}

func TestReadSensorInvalidLinesError(t *testing.T) {
	testSensorID := "foo-sensor"

	ioutilReadFile = func(filename string) ([]byte, error) {
		return []byte("foo"), nil
	}

	actual := readSensor(testSensorID)

	assert.EqualError(t, actual.Error, "sensor data has unexpected number of lines '1'")
	assert.Equal(t, testSensorID, actual.SensorID)
}

func TestReadSensorRegexCompileError(t *testing.T) {
	testErr := "test-read-sensor-regex-compile-error"
	testSensorID := "foo-sensor"

	ioutilReadFile = func(filename string) ([]byte, error) {
		return []byte("foo\nbar"), nil
	}

	regexpCompile = func(expr string) (*regexp.Regexp, error) {
		assert.Equal(t, `t=(\d+)$`, expr)
		return nil, fmt.Errorf(testErr)
	}

	actual := readSensor(testSensorID)

	assert.EqualError(t, actual.Error, fmt.Sprintf("could not compile temperature pattern regular expression 't=(\\d+)$': %s", testErr))
	assert.Equal(t, testSensorID, actual.SensorID)
}

func TestReadSensorInvalidDataError(t *testing.T) {
	testSensorID := "test-sensor"
	ioutilReadFile = func(filename string) ([]byte, error) {
		return []byte("5c 01 4b 46 7f ff 0c 10 d7 : crc=d7 YES\n5c 01 4b 46 7f ff 0c 10 d7\n"), nil
	}

	regexpCompile = regexp.Compile

	actual := readSensor(testSensorID)

	assert.EqualError(t, actual.Error, "temp data did not match expected pattern 't=(\\d+)$'")
	assert.Equal(t, testSensorID, actual.SensorID)
}

func TestReadSensorNonNumericTempError(t *testing.T) {
	testSensorID := "test-sensor"
	ioutilReadFile = func(filename string) ([]byte, error) {
		return []byte("5c 01 4b 46 7f ff 0c 10 d7 : crc=d7 YES\n5c 01 4b 46 7f ff 0c 10 d7 t=21750J\n"), nil
	}

	regexpCompile = func(expr string) (*regexp.Regexp, error) {
		return regexp.Compile(`=(.*)$`)
	}

	actual := readSensor(testSensorID)

	assert.EqualError(t, actual.Error, "could not convert temperature reading '21750J' into integer: strconv.Atoi: parsing \"21750J\": invalid syntax")
	assert.Equal(t, testSensorID, actual.SensorID)
}

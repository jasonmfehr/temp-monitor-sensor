package raspberry_pi

const sensorRootPath = "/sys/bus/w1/devices"
const sensorIDsFile = sensorRootPath + "/w1_bus_master1/w1_master_slaves"
const sensorDataFile = "w1_slave"
const tempPattern = `t=(\d+)$`

type DallasOneWire struct {
}

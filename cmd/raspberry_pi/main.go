package main

import (
	"fmt"

	"github.com/jasonmfehr/temp-monitor-sensor/pkg/awscloud"
	"github.com/jasonmfehr/temp-monitor-sensor/pkg/cloudiface"
	"github.com/jasonmfehr/temp-monitor-sensor/pkg/raspberry_pi"
	"github.com/jasonmfehr/temp-monitor-sensor/pkg/sensoriface"
)

var GitHash = "unknown"
var Version = "unknown"

func main() {
	var sensors sensoriface.Data
	var service cloudiface.Service
	client := cloudiface.Client{
		GitHash: GitHash,
		Version: Version,
	}

	sensors = &raspberry_pi.DallasOneWire{}
	service = &awscloud.Service{}

	if err := service.SendData(client, sensors.GetData()); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

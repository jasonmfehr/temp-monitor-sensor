package cloudiface

import "github.com/jasonmfehr/temp-monitor-sensor/pkg/sensoriface"

type Client struct {
	GitHash string
	Version string
}

type Service interface {
	SendData(client Client, dataPoints []sensoriface.DataPoint) error
}

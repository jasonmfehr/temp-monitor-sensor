package awscloud

import (
	"github.com/jasonmfehr/temp-monitor-sensor/pkg/cloudiface"
	"github.com/jasonmfehr/temp-monitor-sensor/pkg/sensoriface"
)

type Service struct {
}

type serviceData struct {
	ClientInfo cloudiface.Client
	SensorData []sensoriface.DataPoint
}

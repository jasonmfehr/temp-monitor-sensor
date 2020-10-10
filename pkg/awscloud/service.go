package awscloud

import (
	"github.com/jasonmfehr/temp-monitor-sensor/pkg/cloudiface"
	"github.com/jasonmfehr/temp-monitor-sensor/pkg/sensoriface"
)

func (s *Service) SendData(client cloudiface.Client, dataPoints []sensoriface.DataPoint) error {
	// sd := serviceData{
	// 	ClientInfo: client,
	// 	SensorData: dataPoints,
	// }

	return nil
}

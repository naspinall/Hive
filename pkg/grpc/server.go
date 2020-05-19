package grpc

import (
	context "context"

	"github.com/naspinall/Hive/pkg/models"
)

type modelsServer struct {
	ms models.MeasurementService
}

func (s *modelsServer) CreateMeasurement(ctx context.Context, measurement *Measurement) (*Confirmation, error) {
	m := &models.Measurement{
		Value:    measurement.Value,
		DeviceID: int(measurement.DeviceID),
		Type:     measurement.Type,
		Unit:     measurement.Unit,
	}

	err := s.ms.Create(m, ctx)
	if err != nil {
		return nil, err
	}

	return &Confirmation{Reply: 1}, nil

}

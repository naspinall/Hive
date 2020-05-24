package grpc

import (
	context "context"
	"io"

	"google.golang.org/grpc/metadata"

	"github.com/naspinall/Hive/pkg/models"
)

type modelsServer struct {
	ms models.MeasurementService
	us models.UserService
}

func (s *modelsServer) AcceptJWT(token string, ctx context.Context) (context.Context, error) {
	user := &models.User{Token: token}
	return s.us.AcceptToken(user, ctx)
}

// Create a measurement using gRPC
func (s *modelsServer) CreateMeasurement(ctx context.Context, measurement *Measurement) (*Confirmation, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, models.ErrTokenRequired
	}

	authSlice := headers.Get("Authorization")
	// Only expect one Auth header
	if len(authSlice) != 1 {
		return nil, models.ErrBadLogin
	}

	// Validating JWT
	ctx, err := s.AcceptJWT(authSlice[0], ctx)
	if err != nil {
		return nil, err
	}

	m := &models.Measurement{
		Value:    measurement.GetValue(),
		DeviceID: uint(measurement.GetDeviceID()),
		Type:     measurement.GetType(),
		Unit:     measurement.GetUnit(),
	}

	err = s.ms.Create(m, ctx)
	if err != nil {
		return nil, err
	}

	return &Confirmation{Reply: 1}, nil

}

func (s *modelsServer) CreateMeasuremnts(stream MeasurementService_CreateMeasurementsServer) error {

	// Authenticating stream.
	ctx := stream.Context()
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return models.ErrTokenRequired
	}

	authSlice := headers.Get("Authorization")
	// Only expect one Auth header
	if len(authSlice) != 1 {
		return models.ErrBadLogin
	}

	// Validating JWT
	ctx, err := s.AcceptJWT(authSlice[0], ctx)
	if err != nil {
		return err
	}

	// Creating a measurements for each value streamed
	for {
		measurement, err := stream.Recv()
		if err != io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		m := &models.Measurement{
			Value:    measurement.GetValue(),
			DeviceID: uint(measurement.GetDeviceID()),
			Type:     measurement.GetType(),
			Unit:     measurement.GetUnit(),
		}

		err = s.ms.Create(m, ctx)
		if err != nil {
			return err
		}

		return nil
	}
}

func (s *modelsServer) GetMeasurements(stream MeasurementService_GetMeasurementsServer) error {

	// Authenticating stream.
	ctx := stream.Context()
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return models.ErrTokenRequired
	}

	authSlice := headers.Get("Authorization")
	// Only expect one Auth header
	if len(authSlice) != 1 {
		return models.ErrBadLogin
	}

	// Validating JWT
	ctx, err := s.AcceptJWT(authSlice[0], ctx)
	if err != nil {
		return err
	}

	// Creating a measurements for each value streamed
	for {
		// measurement, err := stream.Recv()
		// if err != io.EOF {
		// 	return nil
		// }
		// if err != nil {
		// 	return err
		// }

		// m := &models.Measurement{
		// 	Value:    measurement.GetValue(),
		// 	DeviceID: uint(measurement.GetDeviceID()),
		// 	Type:     measurement.GetType(),
		// 	Unit:     measurement.GetUnit(),
		// }

		// err = s.ms.Create(m, ctx)
		// if err != nil {
		// 	return err
		// }

		// return nil
	}
}

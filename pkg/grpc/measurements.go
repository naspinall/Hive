package grpc

import (
	context "context"
	"errors"
	"io"
	"time"

	"google.golang.org/grpc"

	"github.com/naspinall/Hive/pkg/models"
)

type measurementsHandler struct {
	ms models.MeasurementService
}

func WithMeasurements(server *grpc.Server, ms models.MeasurementService) {
	mh := measurementsHandler{
		ms,
	}
	RegisterMeasurementServiceServer(server, mh)
}

// Create a measurement using gRPC
func (mh measurementsHandler) CreateMeasurement(ctx context.Context, measurement *Measurement) (*Confirmation, error) {

	m := &models.Measurement{
		Value:    measurement.GetValue(),
		DeviceID: uint(measurement.GetDeviceID()),
		Type:     measurement.GetType(),
		Unit:     measurement.GetUnit(),
	}

	err := mh.ms.Create(m, ctx)
	if err != nil {
		return nil, err
	}

	return &Confirmation{Reply: 1}, nil

}

func (mh measurementsHandler) CreateMeasurements(stream MeasurementService_CreateMeasurementsServer) error {

	var buffer []*models.Measurement
	for {
		measurement, err := stream.Recv()
		// Bad read error
		if err != nil {
			return err
		}

		// End of the stream.
		if err == io.EOF {
			//Emptying buffer
			err = mh.ms.CreateMany(buffer, stream.Context())
			if err != nil {
				return err
			}
		}

		// Creating measurement model.
		m := &models.Measurement{
			Value:    measurement.GetValue(),
			DeviceID: uint(measurement.GetDeviceID()),
			Type:     measurement.GetType(),
			Unit:     measurement.GetUnit(),
		}

		// Filling buffer
		if len(buffer) < 100 {
			buffer = append(buffer, m)
		}

		//Mass adding measurements
		err = mh.ms.CreateMany(buffer, stream.Context())
		if err != nil {
			return err
		}
	}

	return nil
}

func (mh measurementsHandler) GetMeasurements(df *DeviceFilter, stream MeasurementService_GetMeasurementsServer) error {

	// Authenticating stream.
	ctx := stream.Context()

	// Parsing dates
	dateFrom, err := time.Parse(time.RFC3339, df.DateFrom)
	if err != nil {
		return errors.New("Bad time")
	}

	dateTo, err := time.Parse(time.RFC3339, df.DateTo)
	if err != nil {
		return errors.New("Bad time")
	}

	// Creating model filter
	f := &models.Filter{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Type:     df.Type,
		Offset:   uint(df.Offset),
		Limit:    uint(df.Limit),
	}

	if f.Limit == 0 {
		f.Limit = 100
	}

	// Adding filter to context
	ctx = context.WithValue(ctx, models.FilterContextKey("Filter"), f)

	// Getting Total Filtered Count
	total, err := mh.ms.Count(ctx)
	if err != nil {
		return err
	}

	count := uint(0)
	// Getting measurements in a stream.
	for count < total {

		f.Offset = count

		measurements, err := mh.ms.GetMany(ctx)
		if err != nil {
			return err
		}
		go func() {
			for _, measurement := range measurements {
				m := &Measurement{
					Type:     measurement.Type,
					Value:    measurement.Value,
					Unit:     measurement.Type,
					DeviceID: int64(measurement.DeviceID),
				}
				// TODO work out how to handle errors
				stream.Send(m)
			}
		}()

		count += 100
	}

	return nil
}

package models

type ErrorUnauthorized string
type ErrorNotFound string
type ErrorBadRequest string

func (e ErrorUnauthorized) Error() string {
	return string(e)
}

func (e ErrorNotFound) Error() string {
	return string(e)
}

func (e ErrorBadRequest) Error() string {
	return string(e)
}

const (
	// Read Required
	ErrDeviceReadRequired        = ErrorUnauthorized("Read Device Access Required")
	ErrMeasurementReadRequired   = ErrorUnauthorized("Read Measurement Access Required")
	ErrAlarmsReadRequired        = ErrorUnauthorized("Read Alarms Access Required")
	ErrUsersReadRequired         = ErrorUnauthorized("Read Users Access Required")
	ErrSubscriptionsReadRequired = ErrorUnauthorized("Read Subscriptions Access Required")

	// Write Required
	ErrDeviceWriteRequired        = ErrorUnauthorized("Device Write Access Required")
	ErrMeasurementWriteRequired   = ErrorUnauthorized("Measurement Write Access Required")
	ErrAlarmsWriteRequired        = ErrorUnauthorized("Alarms Write Access Required")
	ErrUsersWriteRequired         = ErrorUnauthorized("Users Write Access Required")
	ErrSubscriptionsWriteRequired = ErrorUnauthorized("Subscriptions Write Access Required")

	//Update Required
	ErrDeviceUpdateRequired        = ErrorUnauthorized("Device Update Access Required")
	ErrMeasurementUpdateRequired   = ErrorUnauthorized("Measurement Update Access Required")
	ErrAlarmsUpdateRequired        = ErrorUnauthorized("Alarms Update Access Required")
	ErrUsersUpdateRequired         = ErrorUnauthorized("Users Update Access Required")
	ErrSubscriptionsUpdateRequired = ErrorUnauthorized("Subscriptions Update Access Required")

	//Delete Required
	ErrDeviceDeleteRequired        = ErrorUnauthorized("Device Delete Access Required")
	ErrMeasurementDeleteRequired   = ErrorUnauthorized("Measurement Delete Access Required")
	ErrAlarmsDeleteRequired        = ErrorUnauthorized("Alarms Delete Access Required")
	ErrUsersDeleteRequired         = ErrorUnauthorized("Users Delete Access Required")
	ErrSubscriptionsDeleteRequired = ErrorUnauthorized("Subscriptions Delete Access Required")

	// ID Required
	ErrInvalidID = ErrorBadRequest("ID Required")
)

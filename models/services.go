package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Services struct {
	Alarm       AlarmService
	Device      DeviceService
	Measurement MeasurementService
	User        UserService
	db          *gorm.DB
}

func NewServices(connectionString string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &Services{
		Alarm:       NewAlarmService(db),
		Device:      NewDeviceService(db),
		Measurement: NewMeasurementService(db),
		User:        NewUserService(db),
		db:          db,
	}, nil
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Alarm{}, &Measurement{}, &Device{}).Error
}

func (s *Services) DestructiveReset() error {
	if err := s.db.DropTable(&User{}, &Alarm{}, &Measurement{}, &Device{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) Close() error {
	return s.db.Close()
}

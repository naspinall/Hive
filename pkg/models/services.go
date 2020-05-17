package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ServicesConfig func(*Services) error

type Services struct {
	Alarm        AlarmService
	Device       DeviceService
	Measurement  MeasurementService
	User         UserService
	Subscription SubscriptionService
	RBAC         *RBACService
	db           *gorm.DB
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Alarm{}, &Measurement{}, &Device{}, &Subscription{}, &AlarmsRole{}, &UsersRole{}, &MeasurementsRole{}, &DevicesRole{}, &SubscriptionsRole{}).Error
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

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}

		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithAlarms() ServicesConfig {
	return func(s *Services) error {
		s.Alarm = NewAlarmService(s.db, s.Subscription)
		return nil
	}
}
func WithMeasurements() ServicesConfig {
	return func(s *Services) error {
		s.Measurement = NewMeasurementService(s.db, s.Subscription)
		return nil
	}
}
func WithUsers(pepper string, jwtKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, jwtKey)
		return nil
	}
}
func WithDevices() ServicesConfig {
	return func(s *Services) error {
		s.Device = NewDeviceService(s.db)
		return nil
	}
}

func WithSubscriptions() ServicesConfig {
	return func(s *Services) error {
		s.Subscription = NewSubscriptionService(s.db)
		return nil
	}
}

func WithRBAC() ServicesConfig {
	return func(s *Services) error {
		s.RBAC = NewRBACService(s.db)
		return nil
	}
}

package models

import (
	"log"

	"github.com/go-redis/redis"
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
	RBAC         RBACService
	Cache        CacheService
	db           *gorm.DB
	cache        *redis.Client
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
	return s.db.AutoMigrate(&User{}, &Alarm{}, &Measurement{}, &Device{}, &Subscription{}, &Role{}).Error
}

func (s *Services) DestructiveReset() error {
	if err := s.db.DropTable(&User{}, &Alarm{}, &Measurement{}, &Device{}, &Subscription{}, &Role{}).Error; err != nil {
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
		log.Println("Successfully connected to Database")

		s.db = db
		return nil
	}
}

func WithCache(addr string, password string, db int) ServicesConfig {
	return func(s *Services) error {
		cache, err := NewRedisCache(addr, password, db)
		if err != nil {
			return err
		}
		s.Cache = cache
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
		s.Device = NewDeviceService(s.db, s.Cache)
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

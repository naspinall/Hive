package models

import "github.com/jinzhu/gorm"

type Services struct {
	RetainService  RetainService
	SessionService SessionService
	WillService    WillService
	db             *gorm.DB
}

type ServicesConfig func(*Services) error

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
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

func WithRetain() ServicesConfig {
	return func(s *Services) error {
		s.RetainService = NewRetainService(s.db)
		return nil
	}
}
func WithSession() ServicesConfig {
	return func(s *Services) error {
		s.SessionService = NewSessionService(s.db)
		return nil
	}
}
func WithWill() ServicesConfig {
	return func(s *Services) error {
		s.WillService = NewWillService(s.db)
		return nil
	}
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&Will{}, &Session{}, &Retain{}).Error
}

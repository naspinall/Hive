package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Will struct {
	ClientID string         `gorm:"primary_key"`
	QoS      uint           `gorm:"not null`
	Message  postgres.Jsonb `gorm:"not null"`
}

type willGorm struct {
	db *gorm.DB
}

func NewWillService(db *gorm.DB) WillService {
	return &willGorm{
		db,
	}
}

type WillService interface {
}

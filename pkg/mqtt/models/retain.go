package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Retain struct {
	ID      uint           `gorm:"primary_key"`
	Topic   string         `gorm:"not null"`
	QoS     uint           `gorm:"not null`
	Message postgres.Jsonb `gorm:"not null"`
}

type retainGorm struct {
	db *gorm.DB
}

func NewRetainService(db *gorm.DB) RetainService {
	return &retainGorm{
		db,
	}
}

type RetainService interface {
}

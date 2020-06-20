package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Session struct {
	ClientID    string `gorm:"primary_key"`
	Username    string `gom:"not null"`
	LastConnect time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

type sessionGorm struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) SessionService {
	return &sessionGorm{
		db,
	}
}

type SessionService interface {
}

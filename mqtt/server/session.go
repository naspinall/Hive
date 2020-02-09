package server

import (
	"github.com/jinzhu/gorm"
	"github.com/naspinall/Hive/models"
)

type SessionService interface {
	SessionDB
}

type SessionDB interface {
	ByID(id uint) (*Session, error)
	ByDevice(id uint) (*Session, error)
	Create(sesh *Session) error
	Update(sesh *Session) error
	Delete(id uint) error
}

func NewSessionService(db *gorm.DB) SessionService {
	return &sessionGorm{
		db: db,
	}
}

type sessionGorm struct {
	db *gorm.DB
}

type Session struct {
	gorm.Model
	Device    models.Device `gorm:"foreignkey:DeviceId`
	DeviceID  uint
	SessionID string
}

func (sg *sessionGorm) ByID(id uint) (*Session, error) {
	var session Session
	if err := sg.db.Where("id = ?", id).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
func (sg *sessionGorm) ByDevice(id uint) (*Session, error) {
	var session Session
	device := models.Device{Model: gorm.Model{ID: id}}
	if err := sg.db.Model(&session).Related(&device).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
func (sg *sessionGorm) Create(sesh *Session) error {
	if err := sg.db.Create(sesh).Error; err != nil {
		return err
	}
	return nil
}
func (sg *sessionGorm) Update(sesh *Session) error {
	if err := sg.db.Save(sesh).Error; err != nil {
		return err
	}

	return nil
}
func (sg *sessionGorm) Delete(id uint) error {
	session := Session{Model: gorm.Model{ID: id}}
	if err := sg.db.Delete(&session).Error; err != nil {
		return err
	}
	return nil
}

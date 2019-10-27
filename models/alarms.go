package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Alarm struct {
	gorm.Model
	Type     string `gorm:"not null"`
	Status   string `gorm:"not null"`
	Severity string `gorm:"not null"`
	DeviceID int
	Device   Device `json:"-"`
}

type alarmGorm struct {
	db *gorm.DB
}

type AlarmService interface {
	AlarmDB
}

type AlarmDB interface {
	ByID(id uint) (*Alarm, error)
	ByDevice(id uint) ([]Alarm, error)
	Create(alarm *Alarm) error
	Update(alarm *Alarm) error
	Delete(id uint) error
}

func NewAlarmService(db *gorm.DB) AlarmService {
	return &alarmGorm{
		db: db,
	}
}

func (ag *alarmGorm) ByDevice(id uint) ([]Alarm, error) {

	device := Device{Model: gorm.Model{ID: id}}
	alarms := []Alarm{}
	if err := ag.db.Model(&device).Related(&alarms).Error; err != nil {
		return nil, err
	}
	fmt.Print(alarms)
	return alarms, nil
}

func (ag *alarmGorm) ByID(id uint) (*Alarm, error) {
	var alarm Alarm
	if err := ag.db.Where("id = ?", id).First(&alarm).Error; err != nil {
		return nil, err
	}

	return &alarm, nil
}

func (ag *alarmGorm) Create(alarm *Alarm) error {
	return ag.db.Create(alarm).Error
}

func (ag *alarmGorm) Update(alarm *Alarm) error {
	return ag.db.Save(alarm).Error
}
func (ag *alarmGorm) Delete(id uint) error {
	alarm := Alarm{Model: gorm.Model{ID: id}}
	return ag.db.Delete(alarm).Error
}

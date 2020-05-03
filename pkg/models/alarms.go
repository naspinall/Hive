package models

import (
	"context"
	"database/sql"
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
	ByID(id uint, ctx context.Context) (*Alarm, error)
	ByDevice(id uint, ctx context.Context) ([]Alarm, error)
	Create(alarm *Alarm, ctx context.Context) error
	Update(alarm *Alarm, ctx context.Context) error
	Delete(id uint, ctx context.Context) error
	Many(count int, ctx context.Context) ([]*Alarm, error)
}

func NewAlarmService(db *gorm.DB) AlarmService {
	return &alarmGorm{
		db: db,
	}
}

func (ag *alarmGorm) ByDevice(id uint, ctx context.Context) ([]Alarm, error) {

	device := Device{Model: gorm.Model{ID: id}}
	alarms := []Alarm{}
	if err := ag.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Model(&device).Related(&alarms).Error; err != nil {
		return nil, err
	}
	fmt.Print(alarms)
	return alarms, nil
}

func (ag *alarmGorm) ByID(id uint, ctx context.Context) (*Alarm, error) {
	var alarm Alarm
	if err := ag.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Where("id = ?", id).First(&alarm).Error; err != nil {
		return nil, err
	}

	return &alarm, nil
}

func (ag *alarmGorm) Many(count int, ctx context.Context) ([]*Alarm, error) {

	var alarms []*Alarm
	if err := ag.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Limit(count).Find(&alarms).Error; err != nil {
		return nil, err
	}
	return alarms, nil
}

func (ag *alarmGorm) Create(alarm *Alarm, ctx context.Context) error {
	return ag.db.BeginTx(ctx, &sql.TxOptions{}).Create(alarm).Error
}

func (ag *alarmGorm) Update(alarm *Alarm, ctx context.Context) error {
	return ag.db.BeginTx(ctx, &sql.TxOptions{}).Save(alarm).Error
}
func (ag *alarmGorm) Delete(id uint, ctx context.Context) error {
	alarm := Alarm{Model: gorm.Model{ID: id}}
	return ag.db.BeginTx(ctx, &sql.TxOptions{}).Delete(alarm).Error
}

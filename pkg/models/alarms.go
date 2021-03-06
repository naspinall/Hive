package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

type Alarm struct {
	gorm.Model
	Type     string `gorm:"not null"`
	Status   string `gorm:"not null"`
	Severity string `gorm:"not null"`
	DeviceID uint
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

type alarmWebhook struct {
	Subscription SubscriptionService
	AlarmDB
}

type alarmAuthorization struct {
	AlarmDB
}

func NewAlarmService(db *gorm.DB, Subscription SubscriptionService) AlarmService {
	return &alarmAuthorization{
		&alarmWebhook{
			Subscription: Subscription,
			AlarmDB: &alarmGorm{
				db: db,
			},
		},
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
	return ag.db.Create(alarm).Error
}

func (ag *alarmGorm) Update(alarm *Alarm, ctx context.Context) error {
	return ag.db.Save(alarm).Error
}
func (ag *alarmGorm) Delete(id uint, ctx context.Context) error {
	alarm := Alarm{Model: gorm.Model{ID: id}}
	return ag.db.Delete(alarm).Error
}

func (aw *alarmWebhook) Create(alarm *Alarm, ctx context.Context) error {
	err := aw.AlarmDB.Create(alarm, ctx)
	if err != nil {
		return err
	}

	err = aw.Subscription.Webhook(alarm.DeviceID, "CREATE", "ALARM", alarm)
	// Don't want to error for a bad webhook, will just log.
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (aw *alarmWebhook) Update(alarm *Alarm, ctx context.Context) error {
	err := aw.AlarmDB.Update(alarm, ctx)
	if err != nil {
		return err
	}

	err = aw.Subscription.Webhook(alarm.DeviceID, "UPDATE", "ALARM", alarm)
	// Don't want to error for a bad webhook, will just log.
	if err != nil {
		log.Println(err)
	}
	return nil
}
func (aw *alarmWebhook) Delete(id uint, ctx context.Context) error {
	alarm, err := aw.AlarmDB.ByID(id, ctx)
	if err != nil {
		return err
	}
	err = aw.AlarmDB.Delete(id, ctx)
	if err != nil {
		return err
	}

	err = aw.Subscription.Webhook(alarm.DeviceID, "DELETE", "ALARM", alarm)
	// Don't want to error for a bad webhook, will just log.
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (aa alarmAuthorization) ByID(id uint, ctx context.Context) (*Alarm, error) {
	uc, err := ExtractUserClaims(ctx)
	ar := uc.Role.Alarms
	if err != nil || ar < 1 {
		return nil, ErrAlarmsReadRequired
	}
	return aa.AlarmDB.ByID(id, ctx)
}
func (aa alarmAuthorization) ByDevice(id uint, ctx context.Context) ([]Alarm, error) {
	uc, err := ExtractUserClaims(ctx)
	ar := uc.Role.Alarms
	if err != nil || ar < 1 {
		return nil, ErrAlarmsReadRequired
	}
	return aa.AlarmDB.ByDevice(id, ctx)
}
func (aa alarmAuthorization) Create(alarm *Alarm, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	ar := uc.Role.Alarms
	if err != nil || ar < 2 {
		return ErrAlarmsWriteRequired
	}
	return aa.AlarmDB.Create(alarm, ctx)
}
func (aa alarmAuthorization) Update(alarm *Alarm, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	ar := uc.Role.Alarms
	if err != nil || ar < 3 {
		return ErrAlarmsUpdateRequired
	}
	return aa.AlarmDB.Update(alarm, ctx)
}
func (aa alarmAuthorization) Delete(id uint, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	ar := uc.Role.Alarms
	if err != nil || ar < 4 {
		return ErrAlarmsDeleteRequired
	}
	return aa.AlarmDB.Delete(id, ctx)
}
func (aa alarmAuthorization) Many(count int, ctx context.Context) ([]*Alarm, error) {
	uc, err := ExtractUserClaims(ctx)
	ar := uc.Role.Alarms
	if err != nil || ar < 1 {
		return nil, ErrAlarmsReadRequired
	}
	return aa.AlarmDB.Many(count, ctx)
}

package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

type Measurement struct {
	gorm.Model
	Type     string  `gorm:"not null"`
	Value    float64 `gorm:"not null"`
	Unit     string  `gorm:"not null"`
	DeviceID int
	Device   Device `json:"-"`
}

type measurementGorm struct {
	db *gorm.DB
}

type MeasurementService interface {
	MeasurementDB
}

type MeasurementDB interface {
	ByID(id uint, ctx context.Context) (*Measurement, error)
	ByDevice(id uint, ctx context.Context) ([]Measurement, error)
	Create(measurement *Measurement, ctx context.Context) error
	Update(measurement *Measurement, ctx context.Context) error
	Delete(id uint, ctx context.Context) error
}

func NewMeasurementService(db *gorm.DB) MeasurementService {
	return &measurementGorm{
		db: db,
	}
}

func (mg *measurementGorm) ByDevice(id uint, ctx context.Context) ([]Measurement, error) {

	device := Device{Model: gorm.Model{ID: id}}
	measurements := []Measurement{}
	if err := mg.db.BeginTx(ctx, &sql.TxOptions{}).Model(&device).Related(&measurements).Error; err != nil {
		return nil, err
	}
	return measurements, nil
}

func (mg *measurementGorm) ByID(id uint, ctx context.Context) (*Measurement, error) {
	var measurement Measurement
	if err := mg.db.BeginTx(ctx, &sql.TxOptions{}).Where("id = ?", id).First(&measurement).Error; err != nil {
		return nil, err
	}

	return &measurement, nil
}

func (mg *measurementGorm) Create(measurement *Measurement, ctx context.Context) error {
	err := mg.db.BeginTx(ctx, &sql.TxOptions{}).Create(measurement).Error
	if err != nil {
		return err
	}
	mg.Callback(measurement, ctx)
	return nil
}

func (mg *measurementGorm) Update(measurement *Measurement, ctx context.Context) error {
	return mg.db.BeginTx(ctx, &sql.TxOptions{}).Save(measurement).Error
}
func (mg *measurementGorm) Delete(id uint, ctx context.Context) error {
	measurement := Measurement{Model: gorm.Model{ID: id}}
	return mg.db.BeginTx(ctx, &sql.TxOptions{}).Delete(measurement).Error
}

func (mg *measurementGorm) Callback(m *Measurement, ctx context.Context) {
	var subscriptions []*Subscription
	device := Device{Model: gorm.Model{ID: uint(m.DeviceID)}}

	err := mg.db.BeginTx(ctx, &sql.TxOptions{}).Model(&device).Related(&subscriptions).Error
	if err != nil {
		log.Println("Cannot Load Subscriptions")
	}

	for _, subscription := range subscriptions {

		b, err := json.Marshal(m)
		if err != nil {
			log.Println("Invalid Callback Measurement JSON")
		}

		resp, err := http.Post(subscription.Url, "application/json", bytes.NewBuffer(b))
		if err != nil {
			log.Println(err)
		}

		log.Printf("Webhook successful for %s, with status %d", subscription.Url, resp.StatusCode)

	}
}

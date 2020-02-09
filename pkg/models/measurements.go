package models

import (
	"fmt"

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
	ByID(id uint) (*Measurement, error)
	ByDevice(id uint) ([]Measurement, error)
	Create(measurement *Measurement) error
	Update(measurement *Measurement) error
	Delete(id uint) error
}

func NewMeasurementService(db *gorm.DB) MeasurementService {
	return &measurementGorm{
		db: db,
	}
}

func (mg *measurementGorm) ByDevice(id uint) ([]Measurement, error) {

	device := Device{Model: gorm.Model{ID: id}}
	measurements := []Measurement{}
	if err := mg.db.Model(&device).Related(&measurements).Error; err != nil {
		return nil, err
	}
	fmt.Print(measurements)
	return measurements, nil
}

func (mg *measurementGorm) ByID(id uint) (*Measurement, error) {
	var measurement Measurement
	if err := mg.db.Where("id = ?", id).First(&measurement).Error; err != nil {
		return nil, err
	}

	return &measurement, nil
}

func (mg *measurementGorm) Create(measurement *Measurement) error {
	return mg.db.Create(measurement).Error
}

func (mg *measurementGorm) Update(measurement *Measurement) error {
	return mg.db.Save(measurement).Error
}
func (mg *measurementGorm) Delete(id uint) error {
	measurement := Measurement{Model: gorm.Model{ID: id}}
	return mg.db.Delete(measurement).Error
}

package models

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type Device struct {
	gorm.Model
	Name      string  `gorm:"not null;unique_index" json:"name"`
	IMEI      string  `json:"imei"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type deviceGorm struct {
	db *gorm.DB
}

type deviceRabbitMQ struct {
	ch *amqp.Channel
}

type DeviceDB interface {

	//Getters
	ByName(name string) (*Device, error)
	ByID(id uint) (*Device, error)

	//Mutators
	Create(device *Device) error
	Update(device *Device) error
	Delete(id uint) error
	Many(count int) ([]*Device, error)
	SearchByName(name string) ([]*Device, error)
}

type DeviceMultiplexer struct {
	DeviceDB
	DevicePubSub
}

type DevicePubSub interface {

	// Queue actions
	Publish(device *Device) error

	// Queue Operations
	AutoCreateQueue() error
	Close() error
}

type DeviceService interface {
	DeviceDB
}

type deviceService struct {
	DeviceDB
}

func newDeviceRabbitMQ(connectionInfo string) (*deviceRabbitMQ, error) {
	conn, err := amqp.Dial(connectionInfo)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &deviceRabbitMQ{
		ch: ch,
	}, nil
}

func newDeviceGorm(connectionInfo string) (*deviceGorm, error) {

	//Creating database connection
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	//Setting logmode to get more verbose logs from the database layer.
	db.LogMode(true)

	return &deviceGorm{
		db: db,
	}, nil

}

func NewDeviceService(db *gorm.DB) DeviceService {
	return &deviceService{
		&deviceGorm{db: db},
	}
}

//Getters
func (dg *deviceGorm) Many(count int) ([]*Device, error) {
	var devices []*Device

	err := dg.db.Limit(count).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (dg *deviceGorm) ByName(name string) (*Device, error) {

	var device Device
	//Getting Device from database.
	err := dg.db.Where("name = ?", name).First(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (dg *deviceGorm) ByID(id uint) (*Device, error) {
	var device Device
	//Getting Device from database.
	err := dg.db.Where("id = ?", id).First(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (dg *deviceGorm) SearchByName(name string) ([]*Device, error) {
	var devices []*Device
	if err := dg.db.Where("name LIKE ?", "%"+name+"%").Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

//Mutators
func (dg *deviceGorm) Create(device *Device) (err error) {
	err = dg.db.Create(device).Error
	return
}

func (dg *deviceGorm) Update(device *Device) (err error) {
	err = dg.db.Save(device).Error
	return
}

func (dg *deviceGorm) Delete(id uint) (err error) {
	device := Device{Model: gorm.Model{ID: id}}
	err = dg.db.Delete(&device).Error
	return
}

func (dr *deviceRabbitMQ) Publish(device *Device) error {

	deviceBytes, err := json.Marshal(&device)
	if err != nil {
		return err
	}

	return dr.ch.Publish(
		"",
		"device",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        deviceBytes,
		},
	)
}

func (dr *deviceRabbitMQ) Close() error {
	return dr.ch.Close()
}

func (dr *deviceRabbitMQ) AutoCreateQueue() error {

	//Creating the device queue if it does not already exist.
	if _, err := dr.ch.QueueDeclare(
		"device",
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	return nil

}

//Getters
func (dm *DeviceMultiplexer) ByName(name string) (*Device, error) {
	return dm.DeviceDB.ByName(name)
}
func (dm *DeviceMultiplexer) ByID(id uint) (*Device, error) {
	return dm.DeviceDB.ByID(id)
}

//Mutators
func (dm *DeviceMultiplexer) Create(device *Device) error {

	if err := dm.Publish(device); err != nil {
		return err
	}
	return dm.DeviceDB.Create(device)
}
func (dm *DeviceMultiplexer) Update(device *Device) error {

	if err := dm.Publish(device); err != nil {
		return err
	}
	return dm.DeviceDB.Update(device)
}
func (dm *DeviceMultiplexer) Delete(id uint) error {

	device, err := dm.ByID(id)
	if err != nil {
		return nil
	}

	if err := dm.Publish(device); err != nil {
		return err
	}

	return dm.DeviceDB.Delete(id)
}

//Database Operations
func (dm *DeviceMultiplexer) AutoMigrate() error {
	return dm.AutoMigrate()
}
func (dm *DeviceMultiplexer) DestructiveReset() error {
	return dm.DestructiveReset()
}
func (dm *DeviceMultiplexer) Close() error {

	if err := dm.DevicePubSub.Close(); err != nil {
		return err
	}

	return nil
}

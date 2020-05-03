package models

import (
	"context"
	"database/sql"
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
	ByName(name string, ctx context.Context) (*Device, error)
	ByID(id uint, ctx context.Context) (*Device, error)

	//Mutators
	Create(device *Device, ctx context.Context) error
	Update(device *Device, ctx context.Context) error
	Delete(id uint, ctx context.Context) error
	Many(count int, ctx context.Context) ([]*Device, error)
	SearchByName(name string, ctx context.Context) ([]*Device, error)
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
func (dg *deviceGorm) Many(count int, ctx context.Context) ([]*Device, error) {
	var devices []*Device

	err := dg.db.BeginTx(ctx, &sql.TxOptions{}).Limit(count).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (dg *deviceGorm) ByName(name string, ctx context.Context) (*Device, error) {

	var device Device
	//Getting Device from database.
	err := dg.db.BeginTx(ctx, &sql.TxOptions{}).Where("name = ?", name).First(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (dg *deviceGorm) ByID(id uint, ctx context.Context) (*Device, error) {
	var device Device
	//Getting Device from database.
	err := dg.db.BeginTx(ctx, &sql.TxOptions{}).Where("id = ?", id).First(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (dg *deviceGorm) SearchByName(name string, ctx context.Context) ([]*Device, error) {
	var devices []*Device
	if err := dg.db.BeginTx(ctx, &sql.TxOptions{}).Where("name LIKE ?", "%"+name+"%").Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

//Mutators
func (dg *deviceGorm) Create(device *Device, ctx context.Context) (err error) {
	err = dg.db.BeginTx(ctx, &sql.TxOptions{}).Create(device).Error
	return
}

func (dg *deviceGorm) Update(device *Device, ctx context.Context) (err error) {
	err = dg.db.BeginTx(ctx, &sql.TxOptions{}).Save(device).Error
	return
}

func (dg *deviceGorm) Delete(id uint, ctx context.Context) (err error) {
	device := Device{Model: gorm.Model{ID: id}}
	err = dg.db.BeginTx(ctx, &sql.TxOptions{}).Delete(&device).Error
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
func (dm *DeviceMultiplexer) ByName(name string, ctx context.Context) (*Device, error) {
	return dm.DeviceDB.ByName(name, ctx)
}
func (dm *DeviceMultiplexer) ByID(id uint, ctx context.Context) (*Device, error) {
	return dm.DeviceDB.ByID(id, ctx)
}

//Mutators
func (dm *DeviceMultiplexer) Create(device *Device, ctx context.Context) error {

	if err := dm.Publish(device); err != nil {
		return err
	}
	return dm.DeviceDB.Create(device, ctx)
}
func (dm *DeviceMultiplexer) Update(device *Device, ctx context.Context) error {

	if err := dm.Publish(device); err != nil {
		return err
	}
	return dm.DeviceDB.Update(device, ctx)
}
func (dm *DeviceMultiplexer) Delete(id uint, ctx context.Context) error {

	device, err := dm.ByID(id, ctx)
	if err != nil {
		return nil
	}

	if err := dm.Publish(device); err != nil {
		return err
	}

	return dm.DeviceDB.Delete(id, ctx)
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

package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
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

type deviceCache struct {
	cache CacheService
	DeviceDB
}

type deviceAuditLogger struct {
	DeviceDB
}

type deviceAuthorization struct {
	DeviceDB
}

type deviceRabbitMQ struct {
	ch *amqp.Channel
}

type DeviceDB interface {

	//Getters
	ByName(name string, ctx context.Context) (*Device, error)
	ByID(id uint, ctx context.Context) (*Device, error)
	SearchByName(name string, ctx context.Context) ([]*Device, error)
	Many(count int, ctx context.Context) ([]*Device, error)

	//Mutators
	Create(device *Device, ctx context.Context) error
	Update(device *Device, ctx context.Context) error
	Delete(id uint, ctx context.Context) error
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

func NewDeviceService(db *gorm.DB, cache CacheService) DeviceService {
	return &deviceService{
		&deviceAuthorization{
			&deviceAuditLogger{
				&deviceCache{
					cache,
					&deviceGorm{db: db},
				},
			},
		},
	}
}

//Getters
func (dg *deviceGorm) Many(count int, ctx context.Context) ([]*Device, error) {

	var devices []*Device

	err := dg.db.Limit(count).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (dg *deviceGorm) ByName(name string, ctx context.Context) (*Device, error) {

	var device Device
	//Getting Device from database.
	err := dg.db.Where("name = ?", name).First(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (dg *deviceGorm) ByID(id uint, ctx context.Context) (*Device, error) {
	var device Device
	//Getting Device from database.
	err := dg.db.Where("id = ?", id).First(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (dg *deviceGorm) SearchByName(name string, ctx context.Context) ([]*Device, error) {
	var devices []*Device
	if err := dg.db.Where("name LIKE ?", "%"+name+"%").Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

//Mutators
func (dg *deviceGorm) Create(device *Device, ctx context.Context) (err error) {
	err = dg.db.Create(device).Error
	return
}

func (dg *deviceGorm) Update(device *Device, ctx context.Context) (err error) {
	err = dg.db.Save(device).Error
	return
}

func (dg *deviceGorm) Delete(id uint, ctx context.Context) (err error) {
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

//Getters
func (da *deviceAuditLogger) ByName(name string, ctx context.Context) (*Device, error) {
	uc, err := ExtractUserClaims(ctx)
	if err != nil {
		return nil, ErrNoClaims
	}
	LogGet(uc.UserID, "Devices")
	return da.DeviceDB.ByName(name, ctx)
}
func (da *deviceAuditLogger) ByID(id uint, ctx context.Context) (*Device, error) {
	uc, err := ExtractUserClaims(ctx)
	if err != nil {
		return nil, ErrNoClaims
	}
	LogGet(uc.UserID, "Devices")
	return da.DeviceDB.ByID(id, ctx)
}
func (da *deviceAuditLogger) SearchByName(name string, ctx context.Context) ([]*Device, error) {
	uc, err := ExtractUserClaims(ctx)
	if err != nil {
		return nil, ErrNoClaims
	}
	LogGet(uc.UserID, "Devices")
	return da.DeviceDB.SearchByName(name, ctx)
}
func (da *deviceAuditLogger) Many(count int, ctx context.Context) ([]*Device, error) {
	uc, err := ExtractUserClaims(ctx)
	if err != nil {
		return nil, ErrNoClaims
	}
	LogGet(uc.UserID, "Devices")
	return da.DeviceDB.Many(count, ctx)
}

//Mutators
func (da *deviceAuditLogger) Create(device *Device, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	if err != nil {
		return ErrNoClaims
	}
	LogCreate(uc.UserID, "Devices")
	return da.DeviceDB.Create(device, ctx)
}
func (da *deviceAuditLogger) Update(device *Device, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	if err != nil {
		return ErrNoClaims
	}
	LogUpdate(uc.UserID, "Devices")
	return da.DeviceDB.Update(device, ctx)
}
func (da *deviceAuditLogger) Delete(id uint, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	if err != nil {
		return ErrNoClaims
	}
	LogDelete(uc.UserID, "Devices")
	return da.DeviceDB.Delete(id, ctx)
}

func (da deviceAuthorization) ByName(name string, ctx context.Context) (*Device, error) {
	uc, err := ExtractUserClaims(ctx)
	dr := uc.Role.Devices
	if err != nil || dr < 1 {
		return nil, err
	}
	return da.DeviceDB.ByName(name, ctx)
}
func (da deviceAuthorization) ByID(id uint, ctx context.Context) (*Device, error) {
	uc, err := ExtractUserClaims(ctx)
	dr := uc.Role.Devices
	if err != nil || dr < 1 {
		return nil, ErrDeviceReadRequired
	}
	return da.DeviceDB.ByID(id, ctx)
}
func (da deviceAuthorization) SearchByName(name string, ctx context.Context) ([]*Device, error) {
	uc, err := ExtractUserClaims(ctx)
	dr := uc.Role.Devices
	if err != nil || dr < 1 {
		return nil, ErrDeviceReadRequired
	}
	return da.DeviceDB.SearchByName(name, ctx)
}
func (da deviceAuthorization) Many(count int, ctx context.Context) ([]*Device, error) {
	uc, err := ExtractUserClaims(ctx)
	dr := uc.Role.Devices
	if err != nil || dr < 1 {
		return nil, ErrDeviceReadRequired
	}
	return da.DeviceDB.Many(count, ctx)
}

//Mutators
func (da deviceAuthorization) Create(device *Device, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	dr := uc.Role.Devices
	if err != nil || dr < 2 {
		return ErrDeviceWriteRequired
	}
	return da.DeviceDB.Create(device, ctx)
}

func (da deviceAuthorization) Update(device *Device, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	dr := uc.Role.Devices
	if err != nil || dr < 3 {
		return ErrDeviceUpdateRequired
	}
	return da.DeviceDB.Update(device, ctx)
}
func (da deviceAuthorization) Delete(id uint, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	dr := uc.Role.Devices
	if err != nil || dr < 4 {
		return ErrDeviceDeleteRequired
	}
	return da.DeviceDB.Delete(id, ctx)
}

func (dc deviceCache) Create(device *Device, ctx context.Context) error {

	// Creating Device in DB
	err := dc.DeviceDB.Create(device, ctx)
	if err != nil {
		return err
	}

	// Load into cache as likey will be accessed again soon.
	key := fmt.Sprintf("device-%d", device.ID)
	dc.cache.Set(key, device)

	return nil
}

func (dc deviceCache) ByID(id uint, ctx context.Context) (*Device, error) {

	// Getting from the cache
	key := fmt.Sprintf("device-%d", id)
	fmt.Println(key)
	data, err := dc.cache.Get(key)

	if err == redis.Nil {
		log.Printf("No dice")
		// Getting from database and hydrating cache.
		device, err := dc.DeviceDB.ByID(id, ctx)

		// Setting value in the cache as it doesn't exist.
		dc.cache.Set(key, device)
		return device, err
	} else if err != nil {
		return dc.DeviceDB.ByID(id, ctx)
	}

	// Device typecase
	cacheDevice, ok := data.(Device)

	// Using cached value.
	if ok {
		return &cacheDevice, nil
	}

	// Errors with cache just return as normal.
	return dc.DeviceDB.ByID(id, ctx)
}

func (dc deviceCache) Update(device *Device, ctx context.Context) error {

	// Getting from the cache
	key := fmt.Sprintf("device-%d", device.ID)
	err := dc.DeviceDB.Update(device, ctx)
	if err != nil {
		return err
	}
	// Getting updated device.
	ud, err := dc.DeviceDB.ByID(device.ID, ctx)
	if err != nil {
		return err
	}

	// Refreshing value in cache.
	dc.cache.Set(key, ud)
	return nil
}

func (dc deviceCache) Delete(id uint, ctx context.Context) error {
	err := dc.Delete(id, ctx)
	if err != nil {
		return err
	}
	// Getting from the cache
	key := fmt.Sprintf("device-%d", id)
	err = dc.cache.Remove(key)
	if err != nil {
		return err
	}
	return nil
}

package models

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
)

type Subscription struct {
	gorm.Model
	Url      string
	Type     string
	Action   string
	DeviceID uint
	Device   Device `json:"-"`
}

type SubscriptionMessage struct {
	DeviceID uint        `json:"deviceId"`
	Action   string      `json:"action"`
	Type     string      `json:"type"`
	Payload  interface{} `json:"payload"`
}

type subscriptionGorm struct {
	db *gorm.DB
}

type SubscriptionService interface {
	SubscriptionDB
}

type SubscriptionDB interface {
	ByID(id uint) (*Subscription, error)
	ByDevice(id uint) ([]Subscription, error)
	Create(subscription *Subscription) error
	Update(subscription *Subscription) error
	Delete(id uint) error
	Many() ([]*Subscription, error)
	Webhook(deviceID uint, action, Type string, data interface{}) error
}

func NewSubscriptionService(db *gorm.DB) SubscriptionService {
	return &subscriptionGorm{
		db: db,
	}
}

type Webhook func(deviceID uint, action, Type string, data interface{}) error

func (sg *subscriptionGorm) ByDevice(id uint) ([]Subscription, error) {

	device := Device{Model: gorm.Model{ID: id}}
	subscriptions := []Subscription{}
	if err := sg.db.Model(&device).Related(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (sg *subscriptionGorm) ByID(id uint) (*Subscription, error) {
	var subscription Subscription
	if err := sg.db.Where("id = ?", id).First(&subscription).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (sg *subscriptionGorm) Many() ([]*Subscription, error) {
	var subscriptions []*Subscription
	if err := sg.db.Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (sg *subscriptionGorm) Create(subscription *Subscription) error {
	return sg.db.Create(subscription).Error
}

func (sg *subscriptionGorm) Update(subscription *Subscription) error {
	return sg.db.Save(subscription).Error
}
func (sg *subscriptionGorm) Delete(id uint) error {
	subscription := Subscription{Model: gorm.Model{ID: id}}
	return sg.db.Delete(subscription).Error
}

func (sg *subscriptionGorm) Webhook(deviceID uint, action, Type string, data interface{}) error {
	var subscriptions []Subscription
	device := Device{Model: gorm.Model{ID: deviceID}}
	if err := sg.db.Where("action = ?", action).Where("type = ?", Type).Model(&device).Related(&subscriptions).Error; err != nil {
		return err
	}

	// Sending all data for each subscription
	for _, subscription := range subscriptions {
		m := SubscriptionMessage{
			Type:     Type,
			DeviceID: deviceID,
			Action:   action,
			Payload:  data,
		}

		b, err := json.Marshal(&m)
		if err != nil {
			return err
		}

		// Buffer for reading into request
		buff := bytes.NewReader(b)
		_, err = http.Post(subscription.Url, "application/json", buff)
		if err != nil {
			return err
		}
	}

	return nil
}

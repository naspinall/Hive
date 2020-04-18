package models

import (
	"github.com/jinzhu/gorm"
)

type Subscription struct {
	gorm.Model
	Url      string
	Channel  string
	DeviceID int
	Device   Device `json:"-"`
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
}

func NewSubscriptionService(db *gorm.DB) SubscriptionService {
	return &subscriptionGorm{
		db: db,
	}
}

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

func (sg *subscriptionGorm) Webhook(deviceID uint, channel string, data interface{}) error {
	var subscriptions []Subscription
	if err := sg.db.Find(&subscriptions).Where("deviceID = ?", deviceID).Where("channel = ?", channel).Error; err != nil {
		return err
	}
	// for _, subscription := range subscriptions {
	// 	b, err := json.Marshal(data)
	// }
	return nil
}

package models

import (
	"bytes"
	"context"
	"database/sql"
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

type subscriptionAuthorization struct {
	SubscriptionDB
}

type SubscriptionService interface {
	SubscriptionDB
}

type SubscriptionDB interface {
	ByID(id uint, ctx context.Context) (*Subscription, error)
	ByDevice(id uint, ctx context.Context) ([]Subscription, error)
	Create(subscription *Subscription, ctx context.Context) error
	Update(subscription *Subscription, ctx context.Context) error
	Delete(id uint, ctx context.Context) error
	Many(ctx context.Context) ([]*Subscription, error)
	Webhook(deviceID uint, action, Type string, data interface{}) error
}

func NewSubscriptionService(db *gorm.DB) SubscriptionService {
	return &subscriptionGorm{
		db: db,
	}
}

type Webhook func(deviceID uint, action, Type string, data interface{}) error

func (sg *subscriptionGorm) ByDevice(id uint, ctx context.Context) ([]Subscription, error) {

	device := Device{Model: gorm.Model{ID: id}}
	subscriptions := []Subscription{}
	if err := sg.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Model(&device).Related(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (sg *subscriptionGorm) ByID(id uint, ctx context.Context) (*Subscription, error) {
	var subscription Subscription
	if err := sg.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Where("id = ?", id).First(&subscription).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (sg *subscriptionGorm) Many(ctx context.Context) ([]*Subscription, error) {
	var subscriptions []*Subscription
	if err := sg.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (sg *subscriptionGorm) Create(subscription *Subscription, ctx context.Context) error {
	return sg.db.BeginTx(ctx, &sql.TxOptions{}).Create(subscription).Error
}

func (sg *subscriptionGorm) Update(subscription *Subscription, ctx context.Context) error {
	return sg.db.BeginTx(ctx, &sql.TxOptions{}).Save(subscription).Error
}
func (sg *subscriptionGorm) Delete(id uint, ctx context.Context) error {
	subscription := Subscription{Model: gorm.Model{ID: id}}
	return sg.db.BeginTx(ctx, &sql.TxOptions{}).Delete(subscription).Error
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

func (sa subscriptionAuthorization) ByID(id uint, ctx context.Context) (*Subscription, error) {
	uc, err := ExtractUserClaims(ctx)
	sr := uc.Roles.Subscription
	if err != nil || sr < 1 {
		return nil, err
	}
	return sa.SubscriptionDB.ByID(id, ctx)
}
func (sa subscriptionAuthorization) ByDevice(id uint, ctx context.Context) ([]Subscription, error) {
	uc, err := ExtractUserClaims(ctx)
	sr := uc.Roles.Subscription
	if err != nil || sr < 1 {
		return nil, err
	}
	return sa.SubscriptionDB.ByDevice(id, ctx)
}
func (sa subscriptionAuthorization) Create(subscription *Subscription, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	sr := uc.Roles.Subscription
	if err != nil || sr < 2 {
		return err
	}
	return sa.SubscriptionDB.Create(subscription, ctx)
}
func (sa subscriptionAuthorization) Update(subscription *Subscription, ctx context.Context) error {
	return sa.SubscriptionDB.Update(subscription, ctx)
	uc, err := ExtractUserClaims(ctx)
	sr := uc.Roles.Subscription
	if err != nil || sr < 3 {
		return err
	}
	return sa.SubscriptionDB.Update(subscription, ctx)
}
func (sa subscriptionAuthorization) Delete(id uint, ctx context.Context) error {
	uc, err := ExtractUserClaims(ctx)
	sr := uc.Roles.Subscription
	if err != nil || sr < 4 {
		return err
	}
	return sa.SubscriptionDB.Delete(id, ctx)
}
func (sa subscriptionAuthorization) Many(ctx context.Context) ([]*Subscription, error) {
	uc, err := ExtractUserClaims(ctx)
	sr := uc.Roles.Subscription
	if err != nil || sr < 1 {
		return nil, err
	}
	return sa.SubscriptionDB.Many(ctx)
}

package models

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
)

type AccessLevel uint

type AlarmsRole struct {
	gorm.Model
	AccessLevel AccessLevel `gorm:"not null" json:"access"`
	UserID      uint        `gorm:"unique_index`
	User        User        `json:"user"`
}

type UsersRole struct {
	gorm.Model
	AccessLevel AccessLevel `gorm:"not null" json:"access"`
	UserID      uint        `gorm:"unique_index`
	User        User        `json:"user"`
}

type MeasurementsRole struct {
	gorm.Model
	AccessLevel AccessLevel `gorm:"not null" json:"access"`
	UserID      uint        `gorm:"unique_index`
	User        User        `json:"user"`
}

type DevicesRole struct {
	gorm.Model
	AccessLevel AccessLevel `gorm:"not null" json:"access"`
	UserID      uint        `gorm:"unique_index`
	User        User        `json:"user"`
}

type SubscriptionsRole struct {
	gorm.Model
	AccessLevel AccessLevel `gorm:"not null" json:"access"`
	UserID      uint        `gorm:"unique_index`
	User        User        `json:"user"`
}

type alarmRoleGorm struct {
	db *gorm.DB
}

type userRoleGorm struct {
	db *gorm.DB
}

type measurementRoleGorm struct {
	db *gorm.DB
}

type deviceRoleGorm struct {
	db *gorm.DB
}

type subscriptionsGorm struct {
	db *gorm.DB
}

type RBACService struct {
	Alarms        *alarmRoleGorm
	Users         *userRoleGorm
	Measurements  *measurementRoleGorm
	Devices       *deviceRoleGorm
	Subscriptions *subscriptionsGorm
}

func NewRBACService(db *gorm.DB) *RBACService {
	return &RBACService{
		Alarms: &alarmRoleGorm{
			db: db,
		},
		Users: &userRoleGorm{
			db: db,
		},
		Measurements: &measurementRoleGorm{
			db: db,
		},
		Devices: &deviceRoleGorm{
			db: db,
		},
		Subscriptions: &subscriptionsGorm{
			db: db,
		},
	}
}

func (arg *alarmRoleGorm) Assign(role *AlarmsRole, ctx context.Context) error {
	return arg.db.BeginTx(ctx, &sql.TxOptions{}).Create(role).Error
}

func (arg *alarmRoleGorm) ByUserID(userID uint, ctx context.Context) (AccessLevel, error) {
	var alarmRole AlarmsRole
	user := &User{Model: gorm.Model{ID: userID}}
	if err := arg.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Model(&user).Related(&alarmRole).Error; err != nil {
		return AccessLevel(0), err
	}
	return alarmRole.AccessLevel, nil
}

func (urg *userRoleGorm) Assign(role *UsersRole, ctx context.Context) error {
	return urg.db.BeginTx(ctx, &sql.TxOptions{}).Create(role).Error
}

func (urg *userRoleGorm) ByUserID(userID uint, ctx context.Context) (AccessLevel, error) {
	var userRole UsersRole
	user := &User{Model: gorm.Model{ID: userID}}
	if err := urg.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Model(&user).Related(&userRole).Error; err != nil {
		return AccessLevel(0), err
	}
	return userRole.AccessLevel, nil
}

func (mrg *measurementRoleGorm) Assign(role *MeasurementsRole, ctx context.Context) error {
	return mrg.db.BeginTx(ctx, &sql.TxOptions{}).Create(role).Error
}

func (mrg *measurementRoleGorm) ByUserID(userID uint, ctx context.Context) (AccessLevel, error) {
	var measurementRole MeasurementsRole
	user := &User{Model: gorm.Model{ID: userID}}
	if err := mrg.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Model(&user).Related(&measurementRole).Error; err != nil {
		return AccessLevel(0), err
	}
	return measurementRole.AccessLevel, nil
}

func (srg *subscriptionsGorm) Assign(role *SubscriptionsRole, ctx context.Context) error {
	return srg.db.BeginTx(ctx, &sql.TxOptions{}).Create(role).Error
}

func (srg *subscriptionsGorm) ByUserID(userID uint, ctx context.Context) (AccessLevel, error) {
	var subscriptionRole SubscriptionsRole
	user := &User{Model: gorm.Model{ID: userID}}
	if err := srg.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Model(&user).Related(&subscriptionRole).Error; err != nil {
		return AccessLevel(0), err
	}
	return subscriptionRole.AccessLevel, nil
}

func (drg *deviceRoleGorm) Assign(role *DevicesRole, ctx context.Context) error {
	return drg.db.BeginTx(ctx, &sql.TxOptions{}).Create(role).Error
}

func (drg *deviceRoleGorm) ByUserID(userID uint, ctx context.Context) (AccessLevel, error) {
	var deviceRole DevicesRole
	user := &User{Model: gorm.Model{ID: userID}}
	if err := drg.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Model(&user).Related(&deviceRole).Error; err != nil {
		return AccessLevel(0), err
	}
	return deviceRole.AccessLevel, nil
}

package models

import (
	"context"

	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	Alarms        uint `json:"alarms"`
	Users         uint `json:"users"`
	Measurements  uint `json:"measurements"`
	Devices       uint `json:"devices"`
	Subscriptions uint `json:"subscriptions"`
	UserID        uint `json:"userID"`
	User          User `json:"-"`
}

type rbacGorm struct {
	db *gorm.DB
}

type RBACService interface {
	Assign(role *Role, ctx context.Context) error
	ByUserID(id uint, ctx context.Context) (*Role, error)
}

func NewRBACService(db *gorm.DB) RBACService {
	return &rbacGorm{
		db: db,
	}
}

func (rg *rbacGorm) Assign(role *Role, ctx context.Context) error {
	user := User{Model: gorm.Model{ID: role.UserID}}
	if rg.db.Model(&user).Related(&Role{}).RecordNotFound() {
		return rg.db.Create(role).Error
	}
	return rg.db.Model(role).Updates(role).Error
}

func (rg *rbacGorm) ByUserID(id uint, ctx context.Context) (*Role, error) {
	var role Role
	user := User{Model: gorm.Model{ID: id}}
	if err := rg.db.Model(&user).Related(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

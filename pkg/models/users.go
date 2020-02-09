package models

import (
	"github.com/jinzhu/gorm"
)

// User Structs

type User struct {
	gorm.Model
	Username     string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	DisplayName  string `gorm:"not null"`
}

type userGorm struct {
	db *gorm.DB
}

// User Interfaces
// Add Validation here later I guess
type UserService interface {
	UserDB
}

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

type userService struct {
	UserDB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		&userGorm{
			db: db,
		},
	}
}

func newUserGorm(connectionString string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &userGorm{
		db: db,
	}, nil
}

// Implementing the UserDB Interface
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	if err := ug.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	if err := ug.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(user).Error
}

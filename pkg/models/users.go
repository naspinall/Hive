package models

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type modelError string

const pepper = "CMB"

func (e modelError) Error() string {
	return string(e)
}

const (
	ErrNotFound            modelError = "Not found"
	ErrIDInvalid           modelError = "ID provided was invalid"
	ErrPasswordInvalid     modelError = "Invalid password"
	ErrEmailRequired       modelError = "Email required"
	ErrEmailInvalid        modelError = "Email invalid"
	ErrPasswordTooShort    modelError = "Password is too short"
	ErrPasswordRequired    modelError = "Password is required"
	ErrBadLogin            modelError = "Invalid Password or Email"
	ErrDisplayNameRequired modelError = "Display Name Required"
)

// User Structs

type User struct {
	gorm.Model
	Email        string `gorm:"not null;unique_index" json:"email"`
	Password     string `gorm:"not null"  json:"password,omitempty"`
	PasswordHash string `gorm:"not null"  json:"-"`
	DisplayName  string `gorm:"not null"  json:"displayName"`
}

type UserClaims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

type userGorm struct {
	db *gorm.DB
}

// User Interfaces
type UserService interface {
	UserDB
}

type UserDB interface {
	ByID(id uint, ctx context.Context) (*User, error)
	ByEmail(email string, ctx context.Context) (*User, error)

	Create(user *User, ctx context.Context) error
	Update(user *User, ctx context.Context) error
	Delete(id uint, ctx context.Context) error
	Authenticate(email, password string, ctx context.Context) (*User, error)
	Many(ctx context.Context) ([]*User, error)
}

type userValFunc func(*User) error

type userService struct {
	UserDB
}

type userValidator struct {
	UserDB
	emailRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
}

func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db: db}
	return &userService{
		&userValidator{
			UserDB:     ug,
			emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
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
func (ug *userGorm) ByID(id uint, ctx context.Context) (*User, error) {
	var user User
	if err := ug.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func (uv *userValidator) ByID(id uint, ctx context.Context) (*User, error) {
	return uv.UserDB.ByID(id, ctx)
}

func (ug *userGorm) ByEmail(email string, ctx context.Context) (*User, error) {
	var user User
	if err := ug.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (uv *userValidator) ByEmail(email string, ctx context.Context) (*User, error) {
	return uv.UserDB.ByEmail(email, ctx)
}

func (ug *userGorm) Authenticate(email, password string, ctx context.Context) (*User, error) {

	u, err := ug.ByEmail(email, ctx)
	if err != nil {
		return nil, err
	}

	// Adding pepeper to the password.
	toBeCompared := password + pepper

	// Comparing hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(toBeCompared)); err != nil {
		return nil, err
	}
	return u, nil

}

func (uv *userValidator) Authenticate(email, password string, ctx context.Context) (*User, error) {
	user := &User{Email: email, Password: password}
	if err := uv.runUserValFns(user, uv.hasEmail, uv.validEmail, uv.hasPassword); err != nil {
		return nil, err
	}
	return uv.UserDB.Authenticate(email, password, ctx)
}

func (ug *userGorm) Create(user *User, ctx context.Context) error {
	return ug.db.Create(user).Error
}

func (uv *userValidator) Create(user *User, ctx context.Context) error {
	// Ordering in terms of cost.

	if err := uv.runUserValFns(user, uv.hasEmail, uv.validEmail, uv.hasDisplayName, uv.hasPassword, uv.hashPassword, uv.hasPasswordHash); err != nil {
		return err
	}

	return uv.UserDB.Create(user, ctx)
}

func (ug *userGorm) Many(ctx context.Context) ([]*User, error) {
	var users []*User
	if err := ug.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}).Limit(100).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil

}

func (ug *userGorm) Update(user *User, ctx context.Context) error {
	return ug.db.Save(user).Error
}

func (uv *userValidator) Update(user *User, ctx context.Context) error {
	if err := uv.runUserValFns(user, uv.validEmail, uv.validPassword, uv.hashPassword); err != nil {
		return err
	}

	return uv.UserDB.Update(user, ctx)
}

func (ug *userGorm) Delete(id uint, ctx context.Context) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(user).Error
}

func (uv *userValidator) Delete(id uint, ctx context.Context) error {
	return uv.UserDB.Delete(id, ctx)
}

func (uv *userValidator) runUserValFns(u *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(u); err != nil {
			return err
		}
	}

	return nil
}

func (uv *userValidator) hashPassword(user *User) error {
	// Running only if password hasn't already been hashed
	if user.Password == "" {
		return nil
	}

	// Adding pepper to password.
	toBeHashed := user.Password + pepper

	// Hashing using bcrypt, salt is automatically added to the password.
	hash, err := bcrypt.GenerateFromPassword([]byte(toBeHashed), 6)
	if err != nil {
		return err
	}
	// Adding hash to user object
	user.PasswordHash = string(hash)

	// Removing password
	user.Password = ""

	return nil
}

func (uv *userValidator) hasEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) hasPassword(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) hasPasswordHash(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) hasDisplayName(user *User) error {
	if user.DisplayName == "" {
		return ErrDisplayNameRequired
	}
	return nil
}

func (uv *userValidator) validEmail(user *User) error {
	if user.Email == "" {
		return nil
	}
	if uv.emailRegex.Match([]byte(user.Email)) {
		return nil
	}

	return ErrEmailInvalid
}
func (uv *userValidator) validPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) > 8 {
		return nil
	}

	return ErrPasswordInvalid
}

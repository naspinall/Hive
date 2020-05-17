package models

import (
	"context"
	"database/sql"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type modelError string
type userContextKey string

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
	ErrTokenRequired       modelError = "Auth Token Required"
	ErrInvalidClaims       modelError = "Invalid Token Claims"
	ErrInvalidToken        modelError = "Invalid Token"
	ErrNoClaims            modelError = "No Claims"
)

// User Structs

type User struct {
	gorm.Model
	Email        string `gorm:"not null;unique_index" json:"email"`
	Password     string `gorm:"-" json:"password,omitempty"`
	PasswordHash string `gorm:"not null"  json:"-"`
	DisplayName  string `gorm:"not null"  json:"displayName"`
	Token        string `gorm:"-" json:"token,omitempty"`
}

type UserClaims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

type userGorm struct {
	db     *gorm.DB
	pepper string
	jwtKey string
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
	AcceptToken(user *User, ctx context.Context) (context.Context, error)
}

type userValFunc func(*User) error

type userService struct {
	UserDB
	pepper string
}

type userValidator struct {
	UserDB
	emailRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
	tokenRegex    *regexp.Regexp
	pepper        string
}

func NewUserService(db *gorm.DB, pepper, jwtKey string) UserService {
	ug := &userGorm{db: db, pepper: pepper, jwtKey: jwtKey}
	uv := newUserValidator(ug, pepper)
	return &userService{
		UserDB: uv,
	}
}

func newUserValidator(udb UserDB, pepper string) *userValidator {
	return &userValidator{
		UserDB:     udb,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
		pepper:     pepper,
		tokenRegex: regexp.MustCompile(`Bearer ([A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*)`),
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
	toBeCompared := password + ug.pepper

	// Comparing hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(toBeCompared)); err != nil {
		return nil, err
	}

	if err := ug.signToken(u); err != nil {
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
	toBeHashed := user.Password + uv.pepper

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

func (uv *userValidator) hasToken(user *User) error {
	if user.Token == "" {
		return ErrTokenRequired
	}

	return nil
}

func (uv *userValidator) validToken(user *User) error {
	if user.Token == "" {
		return nil
	}

	matches := uv.tokenRegex.FindStringSubmatch(user.Token)
	if len(matches) > 1 {

		user.Token = matches[1]
		return nil
	}

	return ErrTokenRequired

}

func (uv *userValidator) AcceptToken(user *User, ctx context.Context) (context.Context, error) {
	if err := uv.runUserValFns(user, uv.hasToken, uv.validToken); err != nil {
		return ctx, err
	}

	return uv.UserDB.AcceptToken(user, ctx)
}

func (uv *userGorm) AcceptToken(user *User, ctx context.Context) (context.Context, error) {
	uc := &UserClaims{}

	token, err := jwt.ParseWithClaims(
		user.Token,
		uc,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(uv.jwtKey), nil
		},
	)

	if err != nil {
		return ctx, err
	}

	// Setting claims to context TODO do this better, don't need the entire claims in context
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return ctx, ErrInvalidClaims
	}

	claimsContext := context.WithValue(ctx, userContextKey("User"), claims)
	return claimsContext, nil
}

func (ug *userGorm) signToken(user *User) error {
	var err error
	claims := UserClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    "Hive",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims)
	if user.Token, err = token.SignedString([]byte(ug.jwtKey)); err != nil {
		return err
	}

	return nil
}

func ExtractUserClaims(ctx context.Context) (*UserClaims, error) {
	claims, ok := ctx.Value(userContextKey("User")).(*UserClaims)
	if ok {
		return claims, nil
	}
	return nil, ErrInvalidClaims
}

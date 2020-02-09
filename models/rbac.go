package models

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model
	Name        string
	DisplayName string
}

type RoleAssignment struct {
	gorm.Model
	Role   Role
	RoleID int
	User   User
	UserID int
}

type roleGorm struct {
	db *gorm.DB
}

type roleAssignmentGorm struct {
	db *gorm.DB
}

type RoleAssignmentService interface {
	ByUser(user *User) ([]Role, error)
	ByRole(role *Role) ([]Role, error)
	CreateAssignment(user *User, role *Role) (*RoleAssignment, error)

	Create(*RoleAssignment) error
	Update(*RoleAssignment) error
	Delete(*RoleAssignment) error

	Close() error
	AutoMigrate() error
	DestructiveReset() error
}

type RoleService interface {
	ById(id uint) (*Role, error)
	ByName(name string) (*Role, error)

	Create(role *Role)
	Update(role *Role)
	Delete(role *Role)

	Close() error
	AutoMigrate() error
	DestructiveReset() error
}

func (rg *roleGorm) ById(id uint) (*Role, error) {
	var role Role
	if err := rg.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
func (rg *roleGorm) ByName(name string) (*Role, error) {
	var role Role
	if err := rg.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (rg *roleGorm) Create(role *Role) error {
	return rg.db.Create(role).Error
}
func (rg *roleGorm) Update(role *Role) error {
	return rg.db.Save(role).Error

}
func (rg *roleGorm) Delete(id uint) error {
	role := Role{Model: gorm.Model{ID: id}}
	return rg.db.Delete(role).Error
}

func (rg *roleGorm) Close() error {
	return rg.db.Close()
}
func (rg *roleGorm) AutoMigrate() error {
	return rg.db.AutoMigrate(&Role{}).Error
}
func (rg *roleGorm) DestructiveReset() error {
	if err := rg.db.DropTable(&Role{}).Error; err != nil {
		return err
	}

	return rg.AutoMigrate()
}

func (rag *roleAssignmentGorm) ByUser(user *User) ([]RoleAssignment, error) {
	var roles []RoleAssignment
	if err := rag.db.Model(&RoleAssignment{}).Related(user).Find(roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
func (rag *roleAssignmentGorm) ByRole(role *Role) ([]RoleAssignment, error) {
	var roles []RoleAssignment
	if err := rag.db.Model(&role).Related(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
func (rag *roleAssignmentGorm) CreateAssignment(user *User, role *Role) error {
	roleAssignment := RoleAssignment{
		User: *user,
		Role: *role,
	}
	return rag.Create(&roleAssignment)
}
func (rag *roleAssignmentGorm) Create(ra *RoleAssignment) error {
	return rag.db.Create(ra).Error
}
func (rag *roleAssignmentGorm) Update(ra *RoleAssignment) error {
	return rag.db.Save(ra).Error
}
func (rag *roleAssignmentGorm) Delete(id uint) error {
	roleAssignment := RoleAssignment{Model: gorm.Model{ID: id}}
	return rag.db.Delete(roleAssignment).Error
}
func (rag *roleAssignmentGorm) Close() error {
	return rag.db.Close()
}
func (rag *roleAssignmentGorm) AutoMigrate() error {
	return rag.db.AutoMigrate(&RoleAssignment{}).Error
}
func (rag *roleAssignmentGorm) DestructiveReset() error {
	if err := rag.db.DropTable(&RoleAssignment{}).Error; err != nil {
		return err
	}
	return rag.AutoMigrate()
}

package ormgo

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type BaseFromGorm struct {
	gorm.Model
}

type Base struct {
	Uuid      string `gorm:"primary_key;type:varchar(36);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   int `gorm:"type:tinyint(1);default null"`
}

type User struct {
	Base
	DeletedAt *time.Time
	Name      string `gorm:"type:varchar(36);not null"`
	Age       int8   `gorm:"type:tinyint(3)"`
	Sex       string `gorm:"type:enum('men', 'women')"`
}

type Users []User

type Role struct {
	Base
	RoleName string `gorm:"column:rolename;type:varchar(16);unique;not null"`
}

type Roles []Role

type Assignment struct {
	UserUuId string `gorm:"column:user_uuid;primary_key;type:varchar(36)"`
	RoleUuId string `gorm:"column:role_uuid;primary_key;type:varchar(36)"`
}

func (self *User) TableName() string {
	return "users"
}

func (self *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("uuid", uuid.NewV4().String())
	return nil
}

func (self *Role) TableName() string {
	return "roles"
}

func (self *Role) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("uuid", uuid.NewV4().String())
	return nil
}

func (self *Assignment) TableName() string {
	return "assignments"
}

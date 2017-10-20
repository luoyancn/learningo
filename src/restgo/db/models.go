package db

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type BaseFromGorm struct {
	gorm.Model
}

type Base struct {
	Uuid      string    `gorm:"primary_key;type:varchar(36);not null" json:"uuid"`
	CreatedAt time.Time `json:"created_at,string"`
	UpdatedAt time.Time `json:"updated_at,string"`
	Deleted   int       `gorm:"type:tinyint(1);default null" json:"-"`
}

type User struct {
	Base
	DeletedAt *time.Time `json:"deleted_at,string,omitempty"`
	Name      string     `gorm:"type:varchar(36);not null" json:"name"`
	Age       int8       `gorm:"type:tinyint(3)" json:"age"`
	Sex       string     `gorm:"type:enum('men', 'women')" json:"sex"`
}

type Users struct {
	Member []User `json:"users"`
}

type Role struct {
	Base
	RoleName string `gorm:"column:rolename;type:varchar(16);unique;not null" json:"rolename"`
}

type Roles struct {
	Member []Role `json:"roles"`
}

type Assignment struct {
	UserUuId string `gorm:"column:user_uuid;primary_key;type:varchar(36)" json:"user_uuid"`
	RoleUuId string `gorm:"column:role_uuid;primary_key;type:varchar(36)" json:"role_uuid"`
}

type AssignmentUserRole struct {
	Useruuid string `json:"user_uuid"`
	Roleuuid string `json:"role_uuid"`
	Name     string `json:"name"`
	Rolename string `json:"rolename"`
}

type AssignmentUserRoles struct {
	Member []AssignmentUserRole `json:"ass_user_roles"`
}

type Assignments []Assignment

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

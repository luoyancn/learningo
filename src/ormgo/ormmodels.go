package ormgo

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BaseFromGorm struct {
	gorm.Model
}

type Base struct {
	Id        int    `gorm:"auto_increment"`
	Uuid      string `gorm:"primary_key;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Deleted   int `gorm:"type:tinyint(1);default null"`
}

type User struct {
	Base
	Name string `gorm:"type:varchar(36);not null"`
	Age  int    `gorm:"type:int(3)"`
}

func (self User) TableName() string {
	return "users"
}

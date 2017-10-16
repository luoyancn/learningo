package db

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var once sync.Once
var orm_db *gorm.DB

func init() {
	var err error
	once.Do(func() {
		orm_db, err = gorm.Open("mysql",
			"golang:golang@tcp(192.168.137.30:3306)/golang")
		if nil != err {
			fmt.Printf("Cannot init mysql connection:%v\n", err)
			orm_db = nil
		}
		orm_db.LogMode(true)
	})
}

func MigrateDB() {
	orm_db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&User{}, &Role{}, &Assignment{})
	orm_db.Model(&Assignment{}).AddForeignKey(
		"user_uuid", "users(uuid)", "RESTRICT", "RESTRICT")
	orm_db.Model(&Assignment{}).AddForeignKey(
		"role_uuid", "roles(uuid)", "RESTRICT", "RESTRICT")
}

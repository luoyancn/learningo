package ormgo

import (
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
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
	})
}

func MigrateDB() {
	if !orm_db.HasTable(&User{}) {
		orm_db.AutoMigrate(&User{})
	}
}

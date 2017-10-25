package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var once sync.Once
var orm_db *gorm.DB

func InitDbConnection() {
	var err error
	once.Do(func() {
		orm_db, err = gorm.Open(
			"mysql", viper.GetString("database.connection"))
		if nil != err {
			fmt.Printf("Cannot init mysql connection:%v\n", err)
			orm_db = nil
			os.Exit(-2)
		}
		orm_db.LogMode(viper.GetBool("database.debug"))
	})
}

func MigrateDB() {
	InitDbConnection()
	orm_db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&User{}, &Role{}, &Assignment{})
	orm_db.Model(&Assignment{}).AddForeignKey(
		"user_uuid", "users(uuid)", "RESTRICT", "RESTRICT")
	orm_db.Model(&Assignment{}).AddForeignKey(
		"role_uuid", "roles(uuid)", "RESTRICT", "RESTRICT")
}

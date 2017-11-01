package db

import (
	"fastrest/logging"
	"os"
	"sync"
	"time"

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
			logging.TRACE.Fatalf("Cannot init database connection:%v\n", err)
			orm_db = nil
			os.Exit(-2)
		}
		if viper.GetBool("database.debug") {
			orm_db.LogMode(true)
			orm_db.SetLogger(logging.DEBUG)
		}
		orm_db.DB().SetConnMaxLifetime(
			viper.GetDuration("database.max_time") * time.Second)
		orm_db.DB().SetMaxIdleConns(viper.GetInt("database.max_idle"))
		orm_db.DB().SetMaxOpenConns(viper.GetInt("database.max_open"))
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

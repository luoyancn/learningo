package db

import (
	"oceanstack/conf"
	"oceanstack/logging"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var once sync.Once
var orm_db *gorm.DB

func InitDbConnection() {
	var err error
	once.Do(func() {
		orm_db, err = gorm.Open(
			"mysql", conf.DATABASE_CONNECTION)
		if nil != err {
			orm_db = nil
			logging.LOG.Panicf("Cannot init database connection:%v\n", err)
		}
		orm_db.LogMode(conf.DATABASE_DEBUG_MODE)
		orm_db.Set("gorm:table_options", "ENGINE=InnoDB").Set(
			"gorm:table_options", "character set=utf8").Set(
			"gorm:table_options", "collate=utf8_general_ci")
		orm_db.DB().SetConnMaxLifetime(
			conf.DATABASE_MAX_TIME_MIN)
		orm_db.DB().SetMaxIdleConns(conf.DATABASE_MAX_IDLE)
		orm_db.DB().SetMaxOpenConns(conf.DATABASE_MAX_OPEN)
	})
}

func MigrateDB() {
	InitDbConnection()
	orm_db.AutoMigrate(&User{}, &Role{}, &Assignment{})
	orm_db.Model(&Assignment{}).AddForeignKey(
		"user_uuid", "users(uuid)", "RESTRICT", "RESTRICT")
	orm_db.Model(&Assignment{}).AddForeignKey(
		"role_uuid", "roles(uuid)", "RESTRICT", "RESTRICT")
}

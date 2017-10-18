package ormgo

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
			"golang:golang@tcp(192.168.137.30:3306)/golang?parseTime=true")
		if nil != err {
			fmt.Printf("Cannot init mysql connection:%v\n", err)
			orm_db = nil
		}
		//orm_db.LogMode(true)
	})
}

func MigrateDB() {
	if !orm_db.HasTable(&User{}) {
		//orm_db.AutoMigrate(&User{})
		orm_db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
			&User{}, &Role{}, &Assignment{})
		//orm_db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Role{})
		//orm_db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Assignment{})
		orm_db.Model(&Assignment{}).AddForeignKey(
			"user_uuid", "users(uuid)", "RESTRICT", "RESTRICT")
		orm_db.Model(&Assignment{}).AddForeignKey(
			"role_uuid", "roles(uuid)", "RESTRICT", "RESTRICT")
	}

	var query_role Role
	user := User{Name: "zhangjl", Age: 29, Sex: "men"}
	orm_db.Create(&user)

	role := Role{RoleName: "admin"}
	orm_db.Where(&role).Find(&query_role)
	roleuuid := query_role.Uuid
	if query_role.RoleName != role.RoleName {
		orm_db.Create(&role)
		roleuuid = role.Uuid
	}
	orm_db.Create(&Assignment{UserUuId: user.Uuid, RoleUuId: roleuuid})

	var users Users
	var roles Roles
	orm_db.Find(&users)
	for _, user := range users {
		fmt.Printf("%v\n", user)
	}
	orm_db.Find(&roles)
	for _, rol := range roles {
		fmt.Printf("%v\n", rol)
	}
}

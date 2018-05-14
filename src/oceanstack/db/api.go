package db

import (
	"encoding/json"
	"oceanstack/exceptions"
	"oceanstack/logging"

	"github.com/jinzhu/gorm"
)

func UserGet(username string, userpass string) (string, error) {
	var user User
	if err := orm_db.Where(
		"name= ? and password = ?",
		username, userpass).Find(&user).Error; nil != err {
		logging.LOG.Errorf("Error:%v", err)
		switch err {
		case gorm.ErrRecordNotFound:
			return "", exceptions.NewNotFoundException("User cannot be found.")
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
			return "", exceptions.NewSQLException("There were errors in sql")
		case gorm.ErrUnaddressable:
			return "", exceptions.NewConnectionException(err.Error())
		default:
			return "", exceptions.NewException(err.Error())
		}
	}
	user_json, err := json.Marshal(user)
	if nil != err {
		return "", exceptions.NewJsonMarshallException(err.Error())
	}
	return string(user_json), nil
}

func UserCreate(user User) (string, error) {
	if err := orm_db.Create(&user).Error; nil != err {
		switch err {
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
			return "", exceptions.NewSQLException(err.Error())
		case gorm.ErrUnaddressable:
			return "", exceptions.NewConnectionException(err.Error())
		default:
			return "", exceptions.NewException(err.Error())
		}
	}
	logging.LOG.Debugf("The user uuid is %s", user.Uuid)
	return user.Uuid, nil
}

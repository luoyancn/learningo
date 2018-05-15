package db

import (
	"oceanstack/exceptions"
	"oceanstack/logging"

	"github.com/jinzhu/gorm"
)

func UserGet(username string, userpass string) (*User, error) {
	var user User
	if err := orm_db.Where(
		"name= ? and password = ?",
		username, userpass).Find(&user).Error; nil != err {
		logging.LOG.Errorf("Error:%v", err)
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, exceptions.NewNotFoundException("User cannot be found.")
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
			return nil, exceptions.NewSQLException("There were errors in sql")
		case gorm.ErrUnaddressable:
			return nil, exceptions.NewConnectionException(err.Error())
		default:
			return nil, exceptions.NewException(err.Error())
		}
	}
	return &user, nil
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

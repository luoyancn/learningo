package db

import (
	"encoding/json"
	"errors"
	"fastrest/exceptions"

	"github.com/jinzhu/gorm"
)

func UserList() (string, error) {
	var users Users
	if err := orm_db.Find(&users.Member).Error; nil != err {
		return "", err
	}
	user_json, err := json.Marshal(users)
	if nil != err {
		return "", err
	}
	return string(user_json), nil
}

func UserGet(userid string) (string, error) {
	var user User
	if err := orm_db.Where("uuid = ?", userid).Find(&user).Error; nil != err {
		switch err {
		case gorm.ErrRecordNotFound:
			return "", exceptions.NewNotFoundException("User", userid)
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return "", errors.New("ERROR: processing with database")
		default:
			return "", err
		}
	}
	user_json, err := json.Marshal(user)
	if nil != err {
		return "", err
	}
	return string(user_json), nil
}

func UserCreate(user User) error {
	if err := orm_db.Create(&user).Error; nil != err {
		switch err {
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return errors.New("ERROR: processing with database")
		default:
			return err
		}
	}
	return nil
}

func UserUpdate(user User, userid string) error {
	var select_user User
	if err := orm_db.Where("uuid = ?", userid).Find(&select_user).Model(
		&select_user).Update(user).Error; nil != err {
		switch err {
		case gorm.ErrRecordNotFound:
			return exceptions.NewNotFoundException("User", userid)
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return errors.New("ERROR: processing with database")
		default:
			return err
		}
	}
	return nil
}

func UserDelete(userid string) error {
	var assignment Assignment
	if err := orm_db.Where("user_uuid = ?", userid).Find(
		&assignment).Delete(&assignment).Error; nil != err {
		switch err {
		case gorm.ErrRecordNotFound:
			break
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return errors.New("ERROR: processing with database")
		default:
			return err
		}
	}
	if err := orm_db.Where("uuid = ?", userid).Delete(
		&User{}).Error; nil != err {
		switch err {
		case gorm.ErrRecordNotFound:
			return exceptions.NewNotFoundException("User", userid)
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return errors.New("ERROR: processing with database")
		default:
			return err
		}
	}
	return nil
}

func RoleList() (string, error) {
	var roles Roles
	if err := orm_db.Find(&roles.Member).Error; nil != err {
		return "", err
	}
	role_json, err := json.Marshal(roles)
	if nil != err {
		return "", err
	}
	return string(role_json), nil
}

func RoleGet(roleid string) (string, error) {
	var role Role
	if err := orm_db.Where("uuid = ?", roleid).Find(&role).Error; nil != err {
		switch err {
		case gorm.ErrRecordNotFound:
			return "", exceptions.NewNotFoundException("Role", roleid)
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return "", errors.New("ERROR: processing with database")
		default:
			return "", err
		}
	}
	role_json, err := json.Marshal(role)
	if nil != err {
		return "", err
	}
	return string(role_json), nil
}

func RoleCreate(role Role) error {
	if err := orm_db.Create(&role).Error; nil != err {
		switch err {
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return errors.New("ERROR: processing with database")
		default:
			return err
		}
	}
	return nil
}

func RoleDelete(roleid string) error {
	var assignment Assignment
	if err := orm_db.Where("role_uuid = ?", roleid).Find(
		&assignment).Delete(&assignment).Error; nil != err {
		switch err {
		case gorm.ErrRecordNotFound:
			break
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return errors.New("ERROR: processing with database")
		default:
			return err
		}
	}
	var role Role
	if err := orm_db.Where("uuid = ?", roleid).Find(&role).Delete(
		&role).Error; nil != err {
		switch err {
		case gorm.ErrRecordNotFound:
			return exceptions.NewNotFoundException("Role", roleid)
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return errors.New("ERROR: processing with database")
		default:
			return err
		}
	}
	return nil
}

func AssgnmentList(userid string) (string, error) {
	var ass_user_role AssignmentUserRoles
	if err := orm_db.Raw(`select t1.name as name, t2.rolename as rolename,
		t3.user_uuid as useruuid , t3.role_uuid as roleuuid
		from users t1, roles t2, assignments t3 where t1.uuid=t3.user_uuid
		and t2.uuid=t3.role_uuid and t1.uuid = ?`, userid).Scan(
		&ass_user_role.Member).Error; nil != err {
		switch err {
		case gorm.ErrRecordNotFound:
			return "", nil
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return "", errors.New("ERROR: processing with database")
		default:
			return "", err
		}
	}
	ass_json, err := json.Marshal(ass_user_role)
	if nil != err {
		return "", err
	}
	return string(ass_json), nil
}

func AssgnmentCreate(ass Assignment) error {
	if err := orm_db.Create(&ass).Error; nil != err {
		switch err {
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return errors.New("ERROR: processing with database")
		default:
			return err
		}
	}
	return nil
}

func AssgnmentDelete(userid string, roleid string) error {
	var ass Assignment
	if err := orm_db.Where("user_uuid = ? and role_uuid = ?",
		userid, roleid).Delete(&ass).Error; nil != err {
		switch err {
		case gorm.ErrCantStartTransaction:
		case gorm.ErrInvalidSQL:
		case gorm.ErrInvalidTransaction:
		case gorm.ErrUnaddressable:
			return errors.New("ERROR: processing with database")
		default:
			return err
		}
	}
	return nil
}

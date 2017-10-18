package db

import "fmt"

func UserList() {
	var users Users
	orm_db.Find(&users)
	for _, user := range users {
		fmt.Printf("%v\n", user)
	}
}

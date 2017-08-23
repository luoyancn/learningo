package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db_conn *sql.DB
var once sync.Once

func init() {
	var err error
	once.Do(func() {
		db_conn, err = sql.Open(
			"mysql", "nova:nova@tcp(192.168.142.9:3306)/nova")
		if nil != err {
			fmt.Printf("Cannot init mysql connection:%v\n", err)
			db_conn = nil
		}
	})
	db_conn.SetMaxOpenConns(64)
	db_conn.SetMaxIdleConns(32)
	db_conn.SetConnMaxLifetime(60 * time.Second)
}

func main() {
	defer db_conn.Close()
	results, err := db_conn.Query("select uuid from instances")
	if nil != err {
		fmt.Printf("Cannot get any result from dabases :%v\n", err)
		os.Exit(2)
	}

	var uuid_str string
	//var uuids []string
	uuids := []string{}
	defer results.Close()
	for results.Next() {
		err := results.Scan(&uuid_str)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			uuids = append(uuids, uuid_str)
		}
	}
	fmt.Printf("%v\n", uuids)

	stmp, err := db_conn.Prepare(
		"select deleted_at from instances")
	if nil != err {
		fmt.Printf("Cannot prepare sql string :%v\n", err)
		os.Exit(2)
	}
	defer stmp.Close()

	stmp_results, err := stmp.Query()
	if nil != err {
		fmt.Printf("Cannot query with prepare sql string :%v\n", err)
		os.Exit(2)
	}
	defer stmp_results.Close()

	var deleted_at sql.NullString
	time_strs := []string{}
	for stmp_results.Next() {
		err := stmp_results.Scan(&deleted_at)
		if nil != err {
			fmt.Printf("%v\n", err)
		} else {
			if deleted_at.Valid {
				time_strs = append(time_strs, deleted_at.String)
			} else {
				time_strs = append(time_strs, "")
			}
		}
	}
	fmt.Printf("%v\n", time_strs)

	delete_stmp, err := db_conn.Prepare(
		"select uuid, hostname from instances where deleted_at = ?")
	if nil != err {
		fmt.Printf("Cannot query with prepare sql string :%v\n", err)
		os.Exit(2)
	}
	defer delete_stmp.Close()

	res := delete_stmp.QueryRow("2017-05-13 08:56:22")
	var hostname string
	err = res.Scan(&uuid_str, &hostname)
	if nil != err {
		fmt.Printf("Cannot query with prepare sql string :%v\n", err)
		os.Exit(2)
	}
	fmt.Printf("uuid %s:\t%s\n", uuid_str, hostname)

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(time.Now().Format("2006年01月02日 15:04:05"))
}

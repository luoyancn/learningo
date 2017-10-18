package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
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

	ctx, cancle := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancle()

	stmp, err := db_conn.PrepareContext(ctx,
		"select deleted_at from instances")
	if nil != err {
		fmt.Printf("Cannot prepare sql string :%v\n", err)
		os.Exit(2)
	}
	defer stmp.Close()

	stmp_results, err := stmp.QueryContext(ctx)
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
}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime/debug"
)

func loadconfig() string {

	mysqluser := "myuser"
	mysqlpass := "mypassword"
	mysqlurls := "192.168.0.32"
	mysqlport := "3306"
	//mysqldb := "test" // 测试环境
	mysqldb := "zp_mysql" // 生产环境

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", mysqluser, mysqlpass, mysqlurls, mysqlport, mysqldb)
}

func check(e error) {
	if e != nil {
		debug.PrintStack()
		log.Fatal(e)
	}
}

func DBInit() (*sql.DB, error) {

	db, err := sql.Open("mysql", loadconfig())

	if err != nil {
		log.Fatalln("mysql open failed ! check config", err)
	}
	return db, nil
}

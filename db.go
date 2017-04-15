package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DBConnectError error

var DB *sql.DB

func initDB(driver, dsn string) error {
	db, _ := sql.Open(driver, dsn )
	err := db.Ping()
	if err != nil {
		return DBConnectError(err)
	}
	defer db.Close()
	return nil
}

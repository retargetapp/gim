package core

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DBConnectError error

func InitDB(driver, dsn string) (*sql.DB, error) {
	db, _ := sql.Open(driver, dsn )
	err := db.Ping()
	if err != nil {
		return db, DBConnectError(err)
	}
	defer db.Close()
	return db, nil
}

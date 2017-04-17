package core

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	MIGRATIONS_TABLE_NAME  = "gim_migrations"
	MIGRATION_TABLE_SCRIPT = "CREATE TABLE `gim_migrations` (\n" +
		"  `version` int(8) unsigned NOT NULL AUTO_INCREMENT,\n" +
		"  `up` text NOT NULL,\n" +
		"  `down` text NOT NULL,\n" +
		"  PRIMARY KEY (`version`)\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

	ERROR_MIGRATION_TABLE_NOT_EXISTS     = "migration_table_not_exists"
	ERROR_MIGRATION_TABLE_INVALID_SCHEMA = "migration_table_invalid_schema"
)

type MigrationsTableNotExistsError error
type MigrationsTableInvalidError error

func InitDB(driver, dsn string) (*sql.DB, error) {
	db, _ := sql.Open(driver, dsn)
	err := db.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

func CheckMigrationsTable(db *sql.DB) error {
	var t, d string
	err := db.QueryRow("SHOW CREATE TABLE `"+MIGRATIONS_TABLE_NAME+"`").Scan(&t, &d)
	if err != nil {
		return errors.New(ERROR_MIGRATION_TABLE_NOT_EXISTS)
	}
	if d != MIGRATION_TABLE_SCRIPT {
		return errors.New(ERROR_MIGRATION_TABLE_INVALID_SCHEMA)
	}

	return nil
}

func CreateMigrationTable(db *sql.DB) error {
	_, err := db.Exec(MIGRATION_TABLE_SCRIPT)
	return err
}

func LoadDBMigrations(db *sql.DB) (map[uint64]*Migration, error) {
	var m = map[uint64]*Migration{}
	r, err := db.Query("SELECT * FROM " + MIGRATIONS_TABLE_NAME)
	if err != nil {
		return m, err
	}
	defer r.Close()

	for r.Next() {
		var u, d string
		var v uint64
		err := r.Scan(&v, &u, &d)
		if err != nil {
			return m, err
		}
		m[v] = &Migration{
			Version: v,
			Up:      u,
			Down:    d,
		}
	}
	return m, err
}

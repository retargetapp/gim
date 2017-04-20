package core

import (
	"database/sql"

	"reflect"

	"strings"

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
	ERROR_MIGRATION_RECORD_NOT_EXISTS    = "migration_record_not_exists"
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
	rs, err := db.Query("DESC " + MIGRATIONS_TABLE_NAME)
	if err != nil {
		return errors.New(ERROR_MIGRATION_TABLE_NOT_EXISTS)
	}
	defer rs.Close()

	type rt struct {
		Field, Type, Null, Key, Extra string
		Default                       sql.NullString
	}

	var s = make(map[uint]rt)
	var i = uint(0)
	for rs.Next() {
		r := rt{}
		err = rs.Scan(&r.Field, &r.Type, &r.Null, &r.Key, &r.Default, &r.Extra)
		if err != nil {
			return err
		}
		s[i] = r
		i++
	}

	tpl := map[uint]rt{
		0: {
			Field:   "version",
			Type:    "int(8) unsigned",
			Null:    "NO",
			Key:     "PRI",
			Default: sql.NullString{},
			Extra:   "auto_increment",
		},
		1: {
			Field:   "up",
			Type:    "text",
			Null:    "NO",
			Key:     "",
			Default: sql.NullString{},
			Extra:   "",
		},
		2: {
			Field:   "down",
			Type:    "text",
			Null:    "NO",
			Key:     "",
			Default: sql.NullString{},
			Extra:   "",
		},
	}

	if !reflect.DeepEqual(s, tpl) {
		return errors.New(ERROR_MIGRATION_TABLE_INVALID_SCHEMA)
	}

	return nil
}

func CreateMigrationTable(db *sql.DB) error {
	_, err := db.Exec(MIGRATION_TABLE_SCRIPT)
	return err
}

func LoadDBMigrations(db *sql.DB) (map[int64]*Migration, error) {
	var m = map[int64]*Migration{}
	r, err := db.Query("SELECT * FROM " + MIGRATIONS_TABLE_NAME)
	if err != nil {
		return m, err
	}
	defer r.Close()

	for r.Next() {
		var u, d string
		var v int64
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

func LoadDBMigration(db *sql.DB, version string) (*Migration, error) {
	var m = Migration{}
	q := "SELECT * FROM `" + MIGRATIONS_TABLE_NAME + "` WHERE `version` = " + version + ""
	err := db.QueryRow(q).Scan(&m.Version, &m.Up, &m.Down)
	if err != nil {
		if err == sql.ErrNoRows {
			return &m, errors.New(ERROR_MIGRATION_RECORD_NOT_EXISTS)
		}
	}
	return &m, err
}

func ApplyMigration(db *sql.DB, m *Migration) error {
	st, err := db.Prepare("INSERT INTO `" + MIGRATIONS_TABLE_NAME + "`(`version`, `up`, `down`) VALUES(?,?,?);")
	if err != nil {
		return err
	}
	_, err = st.Exec(m.Version, m.Up, m.Down)
	if err != nil {
		return err
	}

	qs := strings.Split(m.Up, ";")
	for _, q := range qs {
		_, err = db.Exec(q)
		if err != nil {
			break
		}
	}

	if err != nil {
		db.Exec("DELETE FROM `"+MIGRATIONS_TABLE_NAME+"` WHERE `version` = ?", m.Version)
		return err
	}
	return nil
}

func RevertMigration(db *sql.DB, m *Migration) error {
	qs := strings.Split(m.Up, ";")
	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}

	db.Exec("DELETE FROM `"+MIGRATIONS_TABLE_NAME+"` WHERE `version` = ?", m.Version)
	return nil
}

package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func configCmd(c *cli.Context) error {
	fmt.Println("Config Gin in current directory")
	cfg, err := core.LoadConfig()

	if err == nil {
		fmt.Println("Gim already configured. Leave fields empty to use previous value:")
	} else {
		if err.Error() == core.ERROR_INVALID_CONFIG_FILE_FORMAT {
			fmt.Println("Config file exists but it has invalid format. Setup new config:")
		}

		cfg = &core.Config{
			Driver: "mysql",
			DSN:    "root:password@tcp(127.0.0.1:3306)/gim",
			Src:    "./migrations",
		}
	}

	var db *sql.DB
	for err = errors.New(""); err != nil; {
		fmt.Printf("Database driver (%s):", cfg.Driver)
		fmt.Scanln(&cfg.Driver)

		fmt.Printf("Data source name (%s):", cfg.DSN)
		fmt.Scanln(&cfg.DSN)

		db, err = core.InitDB(cfg.Driver, cfg.DSN)
		if err != nil {
			fmt.Println("Unable connect to DB: ", err.Error())
		} else {
			defer db.Close()
		}
	}

	err = core.CheckMigrationsTable(db)

	if err == nil {
		fmt.Println("Table `gim_migrations` already exists")
	} else {
		switch err.Error() {
		case core.ERROR_MIGRATION_TABLE_INVALID_SCHEMA:
			fmt.Println("Unable to config database, table `gim_migrations` exists but it's not Gim migrations table")
			// TODO: refactor to property exit way
			os.Exit(1)
		case core.ERROR_MIGRATION_TABLE_NOT_EXISTS:
			err = core.CreateMigrationTable(db)
			if err == nil {
				fmt.Println("Table `gim_migrations` created")
			} else {
				fmt.Println("Unable to create `gim_migrations` table:" + err.Error())
				// TODO: refactor to property exit way
				os.Exit(1)
			}
		}
	}


	for err = errors.New(""); err != nil; {
		var path string
		var f *os.File
		fmt.Printf("Migration files sources directory (%s):", cfg.Src)
		fmt.Scanln(&cfg.Src)

		path, err = filepath.Abs(cfg.Src)
		if err != nil {
			fmt.Println("Unable to check resources directory: " + err.Error())
			continue
		}

		f, err = os.Open(path)
		if err != nil {
			fmt.Println("Unable to check resources directory: " + err.Error())
		} else {
			f.Close()
		}
	}

	err = core.SaveConfig(cfg)

	if err != nil {
		fmt.Println("Unable to save config file: ", core.CONFIG_FILE_NAME, err.Error())
		// TODO: refactor to property exit way
		os.Exit(1)
	}
	fmt.Printf("Config %s updated\n", core.CONFIG_FILE_NAME)
	return nil
}

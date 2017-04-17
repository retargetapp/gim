package main

import (
	"os"

	"fmt"

	"database/sql"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func main() {
	app := cli.NewApp()

	app.Name = "Migrate"
	app.Usage = ""
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Config .gim.yml file and install Git hooks",
			Action: initCmd,
		},
		{
			Name:   "config",
			Usage:  "Config DSN and sql migrations resource directory.",
			Action: configCmd,
		},
		{
			Name:   "install",
			Usage:  "Install Git post-merge and post-checkout hooks for current dir repo.",
			Action: installCmd,
		},
		{
			Name:   "state",
			Usage:  "Show Gim migraions state",
			Action: stateCmd,
		},
		{
			Name:   "sync",
			Usage:  "Sync DB according to migartion resources (sql files)",
			Action: syncCmd,
		},
		{
			Name:   "apply",
			Usage:  "apply <version> – Apply custom migration version",
			Action: applyCmd,
		},
		{
			Name:   "revert",
			Usage:  "apply <version> – Revert custom migration version",
			Action: revertCmd,
		},
		{
			Name:   "up",
			Usage:  "up <version> – Apply all migration versions until defined specified",
			Action: upCmd,
		},
		{
			Name:   "down",
			Usage:  "down <version> – Revert all migration versions until defined specified",
			Action: downCmd,
		},
		{
			Name:   "create",
			Usage:  "create – Create migartion source templates",
			Action: createCmd,
		},
	}

	app.Run(os.Args)
}

func loadConfigHelper() *core.Config {
	cfg, err := core.LoadConfig()
	if err != nil {
		switch err.Error() {
		case core.ERROR_INVALID_CONFIG_FILE_FORMAT:
			fmt.Println("Gim config file `.gim.yml` has invalid format. Run `gim config` to resetup config")
		case core.ERROR_LOAD_CONFIG_FILE:
			fmt.Println("Gim config file `.gim.yml` doesn't exists or is unreadable")
		}
		os.Exit(1)
	}
	return cfg
}

func initDBHelper(cfg *core.Config) *sql.DB {
	db, err := core.InitDB(cfg.Driver, cfg.DSN)
	if err != nil {
		fmt.Println("Unable to connect to DB: " + err.Error())
		os.Exit(1)
	}
	return db
}

/*
func init() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Unable to load config from .gim.yml: %v\n", err)
	}
	log.Printf("Config: %#v\n", cfg)

	err = initDB(cfg.Driver, cfg.DSN)

	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
	}

	m, _ :=loadSources(cfg.Src)
	log.Printf("%#v\n", m[0])
	// Check db
	// – Create DB migrations table

	// Load migrations
	//
}*/

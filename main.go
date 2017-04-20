package main

import (
	"os"

	"database/sql"

	"github.com/urfave/cli"
	"github.com/retargetapp/gim/core"
)

func main() {
	app := cli.NewApp()

	app.Name = "Migrate"
	app.Usage = ""
	app.Version = "0.0.9"
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

func loadConfigHelper() (*core.Config, *cli.ExitError) {
	cfg, err := core.LoadConfig()
	if err != nil {
		switch err.Error() {
		case core.ERROR_INVALID_CONFIG_FILE_FORMAT:
			return cfg, cli.NewExitError("Gim config file `.gim.yml` has invalid format. Run `gim config` to resetup config", 1)
		case core.ERROR_LOAD_CONFIG_FILE:
			return cfg, cli.NewExitError("Gim config file `.gim.yml` doesn't exists or is unreadable", 1)
		}
		return cfg, cli.NewExitError("Load config error", 1)
	}
	return cfg, nil
}

func initDBHelper(cfg *core.Config) (*sql.DB, *cli.ExitError) {
	db, err := core.InitDB(cfg.Driver, cfg.DSN)
	if err != nil {
		return db, cli.NewExitError("Unable to connect to DB: "+err.Error(), 1)
	}
	return db, nil
}

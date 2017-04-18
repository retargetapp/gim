package main

import (
	"fmt"

	"path/filepath"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func applyCmd(c *cli.Context) error {
	fmt.Println("Apply custom migration version")
	if !c.Args().Present() {
		return cli.NewExitError("Migration version undefined. Use `gim apply <version>`", 1)
	}

	v := c.Args().Get(0)
	cfg, cerr := loadConfigHelper()
	if cerr != nil {
		return cerr
	}

	p, err := filepath.Abs(cfg.Src)
	if err != nil {
		return cli.NewExitError("Unable to read sources files from source directory", 1)
	}

	mr, err := core.LoadSrcMigration(p, v)
	if err != nil {
		if rfe, ok := err.(core.ResFileError); ok {
			return cli.NewExitError(rfe.Message(), 1)
		} else if err.Error() == core.ERROR_INVALID_SRC_DIRECTORY {
			return cli.NewExitError("Unable to read sources files from source directory", 1)
		} else {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	db, cerr := initDBHelper(cfg)
	if cerr != nil {
		return cerr
	}

	_, err = core.LoadDBMigration(db, v)
	if err == nil {
		return cli.NewExitError("Migration version `"+v+"` alread applied", 1)
	}

	if err != nil {
		if err.Error() != core.ERROR_MIGRATION_RECORD_NOT_EXISTS {
			return cli.NewExitError("Unable to check current migration state. Error: "+err.Error(), 1)
		}
	}

	fmt.Print("Appling migration with version `" + v + "`...")
	err = core.ApplyMigration(db, mr)
	if err != nil {
		fmt.Println("failed.")
		fmt.Println("Unable to apply migration. Error:" + err.Error())
	}

	fmt.Println("ok.")

	return nil
}

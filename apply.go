package main

import (
	"fmt"

	"strconv"

	"github.com/retargetapp/gim/core"
	"github.com/urfave/cli"
)

func applyCmd(c *cli.Context) error {
	fmt.Println("Apply custom migration version")
	if !c.Args().Present() {
		return cli.NewExitError("Migration version undefined. Use `gim apply <version>`", 1)
	}

	v := c.Args().Get(0)
	vInt, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return cli.NewExitError("Invalid version format. ", 1)
	}

	cfg, cerr := loadConfigHelper()
	if cerr != nil {
		return cerr
	}

	mv, err := core.LoadSrcVersions(cfg.Src)
	if err != nil {
		if rfe, ok := err.(core.ResFileError); ok {
			return cli.NewExitError(rfe.Message(), 1)
		} else if err.Error() == core.ERROR_INVALID_SRC_DIRECTORY {
			return cli.NewExitError("Unable to read sources files from source directory", 1)
		} else {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	var mr *core.Migration
	if fn, ok := mv[vInt]; ok {
		mr, err = core.LoadSrcMigration(cfg.Src, fn, vInt)
		if err != nil {
			if rfe, ok := err.(core.ResFileError); ok {
				return cli.NewExitError(rfe.Message(), 1)
			} else if err.Error() == core.ERROR_INVALID_SRC_DIRECTORY {
				return cli.NewExitError("Unable to read sources files from source directory", 1)
			} else {
				return cli.NewExitError(err.Error(), 1)
			}
		}
	} else {
		return cli.NewExitError("No such version found: "+v, 1)
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

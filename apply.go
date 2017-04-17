package main

import (
	"fmt"

	"os"

	"path/filepath"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func applyCmd(c *cli.Context) error {
	fmt.Println("Apply custom migration version")
	if !c.Args().Present() {
		fmt.Println("Migration version undefined. Use `gim apply <version>`")
		os.Exit(1)
	}

	v := c.Args().Get(0)
	cfg := loadConfigHelper()

	p, err := filepath.Abs(cfg.Src)
	if err != nil {
		fmt.Println("Unable to read sources files from source directory")
		os.Exit(1)
	}

	mr, err := core.LoadSrcMigration(p, v)
	if err != nil {
		if rfe, ok := err.(core.ResFileError); ok {
			fmt.Println(rfe.Message())
		} else if err.Error() == core.ERROR_INVALID_SRC_DIRECTORY {
			fmt.Println("Unable to read sources files from source directory")
		} else {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	db := initDBHelper(cfg)
	_, err = core.LoadDBMigration(db, v)
	if err == nil {
		fmt.Println("Migration version `" + v + "` alread applied")
		os.Exit(0)
	}

	if err != nil {
		if err.Error() != core.ERROR_MIGRATION_RECORD_NOT_EXISTS {
			fmt.Println("Unable to check current migration state. Error: " + err.Error())
			os.Exit(1)
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

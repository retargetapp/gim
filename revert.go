package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func revertCmd(c *cli.Context) error {
	fmt.Println("Revert custom migration version")
	if !c.Args().Present() {
		fmt.Println("Migration version undefined. Use `gim revert <version>`")
		os.Exit(1)
	}

	v := c.Args().Get(0)
	cfg := loadConfigHelper()

	db := initDBHelper(cfg)
	m, err := core.LoadDBMigration(db, v)
	if err != nil {
		if err.Error() == core.ERROR_MIGRATION_RECORD_NOT_EXISTS {
			fmt.Println("Unable to revert migration version `" + v + "`. No such applied version")
		} else {
			fmt.Println("Unable to revert migration version `" + v + ". Error: " + err.Error())
		}
		os.Exit(1)
	}

	fmt.Print("Reverting migration with version `" + v + "`...")
	err = core.RevertMigration(db, m)
	if err != nil {
		fmt.Println("failed.")
		fmt.Println("Unable to revert migration. Error:" + err.Error())
	}

	fmt.Println("ok.")

	return nil
}

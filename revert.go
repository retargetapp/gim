package main

import (
	"fmt"

	"github.com/retargetapp/gim/core"
	"github.com/urfave/cli"
)

func revertCmd(c *cli.Context) error {
	fmt.Println("Revert custom migration version")
	if !c.Args().Present() {
		return cli.NewExitError("Migration version undefined. Use `gim revert <version>`", 1)
	}
	v := c.Args().Get(0)

	cfg, cerr := loadConfigHelper()
	if cerr != nil {
		return cerr
	}

	db, cerr := initDBHelper(cfg)
	if cerr != nil {
		return cerr
	}

	m, err := core.LoadDBMigration(db, v)
	if err != nil {
		if err.Error() == core.ERROR_MIGRATION_RECORD_NOT_EXISTS {
			return cli.NewExitError("Unable to revert migration version `"+v+"`. No such applied version", 1)
		} else {
			return cli.NewExitError("Unable to revert migration version `"+v+". Error: "+err.Error(), 1)
		}
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

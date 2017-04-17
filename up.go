package main

import (
	"fmt"

	"strconv"

	"sort"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func upCmd(c *cli.Context) error {
	var vi int64
	var err error

	fmt.Println("Up to custom migration version")
	if !c.Args().Present() {
		return cli.NewExitError("Migration version undefined. Use `gim up <version>`", 1)
	}

	v := c.Args().Get(0)

	if v == "." {
		vi = -1
	} else {
		vi, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return cli.NewExitError("Invalid version format. Error: "+err.Error(), 1)
		}
	}

	cfg := loadConfigHelper()
	db := initDBHelper(cfg)

	md, err := core.LoadDBMigrations(db)
	if err != nil {
		fmt.Println("Unable to load applied migrations from DB: " + err.Error())
	}

	mr, err := core.LoadSrcMigrations(cfg.Src)
	if err != nil {
		if rfe, ok := err.(core.ResFileError); ok {
			return cli.NewExitError(rfe.Message(), 1)
		} else if err.Error() == core.ERROR_INVALID_SRC_DIRECTORY {
			return cli.NewExitError("Unable to read sources files from source directory", 1)
		} else {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	var vs = []int64{}
	for v, _ := range mr {
		vs = append(vs, v)
	}
	sort.Sort(version(vs))

	for _, v := range vs {
		if _, ok := md[v]; ok {
			continue
		}
		if v == vi {
			break
		}
		fmt.Print("Applying migration with version `" + strconv.FormatInt(v, 10) + "`...")
		err = core.ApplyMigration(db, mr[v])
		if err != nil {
			fmt.Println("failed.")
			return cli.NewExitError("Unable to apply migration. Error:"+err.Error(), 1)
		}
		fmt.Println("ok.")
	}

	return nil
}

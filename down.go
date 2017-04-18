package main

import (
	"fmt"

	"strconv"

	"sort"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func downCmd(c *cli.Context) error {
	var vi int64
	var err error

	fmt.Println("Down to custom migration version")
	if !c.Args().Present() {
		return cli.NewExitError("Migration version undefined. Use `gim down <version>`", 1)
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

	cfg, cerr := loadConfigHelper()
	if cerr != nil {
		return cerr
	}

	db, cerr := initDBHelper(cfg)
	if cerr != nil {
		return cerr
	}

	md, err := core.LoadDBMigrations(db)
	if err != nil {
		return cli.NewExitError("Unable to load applied migration from DB: "+err.Error(), 1)
	}

	if _, ok := md[vi]; !ok && vi > 0 {
		return cli.NewExitError("Unable revert down to migration version `"+v+"`. No such applied version", 1)
	}

	var vs = []int64{}
	for v, _ := range md {
		vs = append(vs, v)
	}
	sort.Sort(sort.Reverse(version(vs)))

	for _, v := range vs {
		if v == vi {
			break
		}
		fmt.Print("Reverting migration with version `" + strconv.FormatInt(v, 10) + "`...")
		err = core.RevertMigration(db, md[v])
		if err != nil {
			fmt.Println("failed.")
			return cli.NewExitError("Unable to revert migration. Error:"+err.Error(), 1)
		}
		fmt.Println("ok.")
	}

	return nil
}

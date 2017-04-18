package main

import (
	"fmt"

	"sort"

	"strconv"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func syncCmd(c *cli.Context) error {
	fmt.Println("Sync Gim state")

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
		fmt.Println("Unable to load applied migration from DB: " + err.Error())
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
	for v, _ := range md {
		vs = append(vs, v)
	}
	sort.Sort(sort.Reverse(version(vs)))

	for _, v := range vs {
		if _, ok := mr[v]; ok {
			continue
		}
		fmt.Print("Reverting migration with version `" + strconv.FormatInt(v, 10) + "`...")
		err = core.RevertMigration(db, md[v])
		if err != nil {
			fmt.Println("failed.")
			return cli.NewExitError("Unable to revert migration. Error:"+err.Error(), 1)
		}
		fmt.Println("ok.")
	}

	vs = []int64{}
	for v, _ := range mr {
		vs = append(vs, v)
	}
	sort.Sort(version(vs))

	for _, v := range vs {
		if _, ok := md[v]; ok {
			continue
		}
		fmt.Print("Applying migration with version `" + strconv.FormatInt(int64(v), 10) + "`...")
		err = core.ApplyMigration(db, mr[v])
		if err != nil {
			fmt.Println("failed.")
			return cli.NewExitError("Unable to apply migration. Error:"+err.Error(), 1)
		}
		fmt.Println("ok.")
	}

	return nil
}

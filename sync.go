package main

import (
	"fmt"

	"os"

	"sort"

	"strconv"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func syncCmd(c *cli.Context) error {
	fmt.Println("Sync Gim state")

	cfg := loadConfigHelper()
	db := initDBHelper(cfg)

	md, err := core.LoadDBMigrations(db)
	if err != nil {
		fmt.Println("Unable to load applied migration from DB: " + err.Error())
	}

	mr, err := core.LoadSrcMigrations(cfg.Src)
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

	var vs = []uint64{}
	for v, _ := range md {
		vs = append(vs, v)
	}
	sort.Sort(sort.Reverse(version(vs)))

	for _, v := range vs {
		if _, ok := mr[v]; ok {
			continue
		}
		fmt.Print("Reverting migration with version `" + strconv.FormatInt(int64(v), 10) + "`...")
		err = core.RevertMigration(db, md[v])
		if err != nil {
			fmt.Println("failed.")
			fmt.Println("Unable to revert migration. Error:" + err.Error())
			os.Exit(1)
		}
		fmt.Println("ok.")
	}

	vs = []uint64{}
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
			fmt.Println("Unable to apply migration. Error:" + err.Error())
			os.Exit(1)
		}
		fmt.Println("ok.")
	}

	return nil
}

package main

import (
	"fmt"
	"os"

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
		fmt.Println("Migration version undefined. Use `gim up <version>`")
		os.Exit(1)
	}

	v := c.Args().Get(0)

	if v == "." {
		vi = -1
	} else {
		vi, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Printf("Invalid version format. Error: " + err.Error())
			os.Exit(1)
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
			fmt.Println(rfe.Message())
		} else if err.Error() == core.ERROR_INVALID_SRC_DIRECTORY {
			fmt.Println("Unable to read sources files from source directory")
		} else {
			fmt.Println(err.Error())
		}
		os.Exit(1)
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
			fmt.Println("Unable to apply migration. Error:" + err.Error())
			os.Exit(1)
		}
		fmt.Println("ok.")
	}

	return nil
}

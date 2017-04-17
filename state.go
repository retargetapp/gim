package main

import (
	"fmt"

	"sort"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

type version []int64

func (v version) Len() int           { return len(v) }
func (v version) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v version) Less(i, j int) bool { return v[i] < v[j] }

func stateCmd(c *cli.Context) error {
	fmt.Println("Gim migration version state:")
	cfg := loadConfigHelper()
	db := initDBHelper(cfg)

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

	var vm = make(map[int64]struct{})
	for v, _ := range md {
		vm[v] = struct{}{}
	}
	for v, _ := range mr {
		vm[v] = struct{}{}
	}

	var vs = []int64{}
	for v, _ := range vm {
		vs = append(vs, v)
	}
	sort.Sort(version(vs))

	if len(vs) > 0 {
		fmt.Println("DB version\tResources version\tSync status")
	} else {
		fmt.Println("No migration to checkt state")
		return nil
	}

	for _, v := range vs {
		var n int8
		if _, ok := md[v]; ok {
			fmt.Print(v)
			n--
		} else {
			fmt.Print("    –    ")
		}
		fmt.Print("\t")

		if _, ok := mr[v]; ok {
			fmt.Print(v, "       ")
			n++
		} else {
			fmt.Print("    -           ")
		}
		fmt.Print("\t")
		if n > 0 {
			fmt.Print("Apply while sync")
		} else if n < 0 {
			fmt.Print("Revert while sync")
		} else {
			fmt.Print("Ok")
		}

		fmt.Print("\n")
	}

	return nil
}

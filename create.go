package main

import (
	"fmt"
	"strconv"
	"time"

	"os"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

func createCmd(c *cli.Context) error {
	cfg := loadConfigHelper()
	v := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Print("Create migration source templates...")
	err := core.CreateSrcVersionTpl(cfg.Src, v)
	if err != nil {
		fmt.Println("failed.")
		fmt.Println("Unable to create source templates. Error:" + err.Error())
		os.Exit(1)
	}
	fmt.Println("ok.")
	return nil
}

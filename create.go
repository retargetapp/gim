package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli"
	"github.com/retargetapp/gim/core"
)

func createCmd(c *cli.Context) error {
	cfg, cerr := loadConfigHelper()
	if cerr != nil {
		return cerr
	}

	v := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Print("Create migration source templates...")
	err := core.CreateSrcVersionTpl(cfg.Src, v)
	if err != nil {
		fmt.Println("failed.")
		return cli.NewExitError("Unable to create source templates. Error:"+err.Error(), 1)
	}
	fmt.Println("ok.")
	return nil
}

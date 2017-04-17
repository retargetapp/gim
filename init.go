package main

import (
	"github.com/urfave/cli"
)

func initCmd(c *cli.Context) error {
	configCmd(c)
	installCmd(c)
	return nil
}

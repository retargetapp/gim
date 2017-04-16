package main

import (
	"log"

	"github.com/urfave/cli"
)

func initCmd(c *cli.Context) error {
	log.Println("Init command call")
	configCmd(c)
	installCmd(c)
	return nil
}
